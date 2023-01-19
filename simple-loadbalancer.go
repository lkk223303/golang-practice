/*
* ----------------------------------------------------------------------------
* Copyright (c) 2022-present BigObject Inc.
* All Rights Reserved.
*
* Use of, copying, modifications to, and distribution of this software
* and its documentation without BigObject's written permission can
* result in the violation of U.S., Taiwan and China Copyright and Patent laws.
* Violators will be prosecuted to the highest extent of the applicable laws.
*
* BIGOBJECT MAKES NO REPRESENTATIONS OR WARRANTIES ABOUT THE SUITABILITY OF
* THE SOFTWARE, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED
* TO THE IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
* PARTICULAR PURPOSE, OR NON-INFRINGEMENT.
*
*
* simple-loadbalancer.go
*
* @author:   Grace Chen, Kent Huang
* ----------------------------------------------------------------------------
*/
	
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Peer interface{}

type Factor interface {
	Factor() string
}
type BalancerLite interface {
	Next(factor Factor) (next Peer, c Constrainable)
}

type Balancer interface {
	BalancerLite
}

type FactorComparable interface {
	Factor
	ConstrainedBy(constraint interface{}) (peer Peer, c Constrainable)
}

type FactorString string

func (f FactorString) Factor() string {
	return string(f)
}

const DummyFactor FactorString = ""

type Constrainable interface {
	CanConstrain(o interface{}) (yes bool)
	Check(o interface{}) (satisfied bool)
	Peer
}

// Random algorithm

type randomS struct {
	peers []Peer
	count int64
}

func (s *randomS) Next(factor Factor) (next Peer, c Constrainable) {

	l := int64(len(s.peers))
	ni := atomic.AddInt64(&s.count, inRange(0, l)) % l

	next = s.peers[ni]
	return

}
func randomMain() {
	lb := &randomS{
		peers: []Peer{
			exP("172.16.0.7:3500"),
			exP("172.16.0.8:3500"),
			exP("172.16.0.9:3500"),
			exP("172.16.0.10:3500"),
			exP("172.16.0.11:3500"),
			exP("172.16.0.12:3500"),
			exP("172.16.0.13:3500"),
		},
		count: 0}

	sum := make(map[Peer]int)
	for i := 0; i < 30000000; i++ {
		p, _ := lb.Next(DummyFactor)
		sum[p]++
	}

	for k, v := range sum {
		fmt.Printf("%v: %v\n", k, v)
	}
}

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var seedmu sync.Mutex

func inRange(min, max int64) int64 {
	seedmu.Lock()
	defer seedmu.Unlock()
	return seededRand.Int63n(max-min) + min
}

type exP string

func (s exP) String() string { return string(s) }
