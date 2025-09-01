package main

import (
	"reflect"
	"sync"
	"time"
)

type AttrVal any
type AttrType = reflect.Kind

type store struct {
	keys             map[string]*ValueObject
	maxEvictWaitTime time.Duration // this is the maximum duration to wait when the evictPool is full
	evictPool        chan struct{} // this forces eviction work to be limited to avoid spawning huge number of goroutines --> Try to use a limited worker goroutine pool to evict keys. If workers are all busy for too long, just evict in the current goroutine.
	mu               sync.RWMutex
}

type ValueObject struct {
	attributes map[string]AttrVal  // key is attribute key which is always string, value is the attribute value
	types      map[string]AttrType // key is attribute key which is always string, value is the attribute type
	expiry     time.Duration
}
