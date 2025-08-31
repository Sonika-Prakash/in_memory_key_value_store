package main

import (
	"reflect"
	"sync"
)

type AttrVal any
type AttrType = reflect.Kind

type store struct {
	keys map[string]*ValueObject
	mu   sync.RWMutex
}

type ValueObject struct {
	attributes map[string]AttrVal  // key is attribute key which is always string, value is the attribute value
	types      map[string]AttrType // key is attribute key which is always string, value is the attribute type
}
