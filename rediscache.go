package main

import (
	"fmt"

	redigo "github.com/garyburd/redigo/redis"
)

// redisCache implements the interface of CacheStore using redis.
type redisCache struct {
	pool *RedisWaitPool // Redis connection pool
}

// NewRedisCache news and returns an instance of redisCache
func NewRedisCache(pool *RedisWaitPool) *redisCache {
	return &redisCache{
		pool: pool,
	}
}

// CheckHealth checks the health of redis
func (r *redisCache) CheckHealth() error {
	_, err := r.do("EXISTS", "foo")
	if err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}
	return nil
}

// Set binds the key with the value to the redis.
func (r *redisCache) Set(key string, value string) error {
	reply, err := r.do("SET", key, value)
	if err != nil || reply != "OK" {
		return fmt.Errorf("failed to do Redis SET command: key: %s, val: %s", key, value)
	}
	return nil
}

// Get gets the data from the multiple keys given.
func (r *redisCache) Get(key string) (string, error) {
	reply, err := r.do("GET", key)
	strReply, err := redigo.String(reply, err)
	if err != nil {
		return "", fmt.Errorf("failed to do redis GET command: key: %s", key)
	}
	return strReply, nil
}

// Mget gets the data from the multiple keys given.
func (r *redisCache) Mget(keys []string) ([]string, error) {
	ikeys := make([]interface{}, 0, len(keys))
	for _, k := range keys {
		ikeys = append(ikeys, k)
	}
	reply, err := r.do("MGET", ikeys...)
	arrReply, err := redigo.Strings(reply, err)
	if err != nil {
		return nil, fmt.Errorf("failed to do redis MGET command: keys:%v", keys)
	}
	return arrReply, nil
}

// Del deletes the key from cachestore
func (r *redisCache) Del(key string) error {
	_, err := r.do("DEL", key)
	if err != nil {
		return fmt.Errorf("failed to do redis DEL comamnd: key:%s", key)
	}
	return nil
}

// GetKeys lists the keys with regex
// For example: keys assoc:* will list all keys prefixed with "assoc:"
func (r *redisCache) GetKeys(regexKey string) ([]string, error) {
	arrReply, err := redigo.Strings(r.do("KEYS", regexKey))
	if err != nil {
		return nil, fmt.Errorf("failed to do redis MGET command: keys:%v", regexKey)
	}
	return arrReply, nil
}

// do write the command to the redis.Conn
func (r *redisCache) do(command string, args ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()
	return conn.Do(command, args...)
}
