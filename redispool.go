package main

import (
	"log"
	"sync"
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

// RedisWaitPool implements conditional waiting mechanism for connection pool.
// If the current connection is unavailable, it waits forever until someone else has signaled.
type RedisWaitPool struct {
	*redigo.Pool
	// RedisWaitPool inherits redigo.Pool

	cond *sync.Cond
}

// NewRedisTCPPool wraps RedisWaitPool with tcp connection
func NewRedisTCPPool(idleConn, activeConn int, addr string) *RedisWaitPool {
	return NewRedisWaitPool(idleConn, activeConn,
		func() (redigo.Conn, error) { return redigo.Dial("tcp", addr) })
}

// NewRedisWaitPool allocates a RedisWaitPool instance.
func NewRedisWaitPool(idleConn, activeConn int, dial func() (redigo.Conn, error)) *RedisWaitPool {
	pool := &redigo.Pool{
		MaxIdle:   idleConn,
		MaxActive: activeConn,
		Dial:      dial,
		// NOTE: on TestOnBorrow, we perform a "EXISTS" operation on a dummy key to
		// health check a connection. The existence of the key does not affect the
		// result. As long as the connection/redis is good, no err will be returned.
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("EXISTS", "foo")
			return err
		},
	}

	log.Printf("RedisWaitPool: create a pool(%p) with idle (%d) and active (%d)", pool, idleConn, activeConn)
	return &RedisWaitPool{
		Pool: pool,
		cond: sync.NewCond(&sync.RWMutex{}),
	}
}

// Get gets the waitPoolConn instance.
func (wp *RedisWaitPool) Get() redigo.Conn {
	return &waitPooledConn{
		p:    wp,
		conn: wp.Pool.Get(),
	}
}

// waitPooledConnection implements the interface of Conn and wrap a redigo.Conn inside.
type waitPooledConn struct {
	conn redigo.Conn
	p    *RedisWaitPool
}

// Close closes the pooled connection.
func (c *waitPooledConn) Close() (err error) {
	return c.conn.Close()
}

// Err returns the error of the connection.
func (c *waitPooledConn) Err() error {
	return c.conn.Err()
}

// redisCmd is the Redis command to be retried
type redisCmd func() (interface{}, error)

// retryCmd executes and retries a Redis command. It absorbs
// redigo.ErrPoolExhausted error and treats it as retry signal. Returns error
// otherwise.
func (c *waitPooledConn) retryCmd(cmd redisCmd) (reply interface{}, err error) {
	var isCont bool

	// resetConn resets the conn data.
	// since redigo will catch the status of the conn.
	// we free the current one and make a new conn to keep the conn clean.
	resetConn := func() {
		c.conn = nil
		c.conn = c.p.Pool.Get()
	}

	// tryCmd tries to execute the command on redis.
	// if the catched error is NOSCRIPT and other errors,
	// we just return it and leave redigo.script to handle it.
	// if the catched error is PoolExhausted, we go to wait and try again.
	tryCmd := func() (bool, interface{}, error) {
		reply, err := cmd()
		if err != nil {
			// Retry in case of ErrPoolExhausted.
			return (err == redigo.ErrPoolExhausted), nil, err
		}
		// normal terminal
		return false, reply, nil
	}

	ret := func(reply interface{}, err error) (interface{}, error) {
		c.p.cond.Signal()
		return reply, err
	}

	isCont, reply, err = tryCmd()
	// we try to do the query to avoid the busy locking.
	if !isCont {
		return ret(reply, err)
	}

	// Otherwise, we go to wait-and-retry policy.
	c.p.cond.L.Lock()
	defer c.p.cond.L.Unlock()

	for {
		c.p.cond.Wait()
		resetConn()
		isCont, reply, err = tryCmd()
		if !isCont {
			return ret(reply, err)
		}
	}
}

// Do executes and retries a command to Redis. It retries when receiving
// redigo.ErrPoolExhausted error and returns otherwise.
func (c *waitPooledConn) Do(commandName string, args ...interface{}) (reply interface{},
	err error) {

	cmd := func() (interface{}, error) {
		reply, err = c.conn.Do(commandName, args...)
		return reply, err
	}
	return c.retryCmd(cmd)
}

// Send sends and retries command to Redis.
// NOTE: we do not implement the wait-and-retry in this command.
func (c *waitPooledConn) Send(commandName string, args ...interface{}) error {
	return c.conn.Send(commandName, args...)
}

// Flush flushs the commands which are sent to redis. It assumes caller already
// holds a valid connection.
func (c *waitPooledConn) Flush() error {
	return c.conn.Flush()
}

// Receive receives the response value from the side of redis. It assumes caller
// already holds a valid connection.
func (c *waitPooledConn) Receive() (reply interface{}, err error) {
	return c.conn.Receive()
}
