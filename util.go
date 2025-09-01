package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

func newStore() *store {
	return &store{
		keys:             make(map[string]*ValueObject),
		mu:               sync.RWMutex{},
		maxEvictWaitTime: 2 * time.Second,
		evictPool:        make(chan struct{}, 5),
	}
}

func (s *store) scheduleEviction(key string) {
	select {
	case s.evictPool <- struct{}{}:
		// try to grab a slot for eviction via goroutine
	case <-time.After(s.maxEvictWaitTime):
		// max time to wait to grab a slot which is the key's expiry duration
		// if no slot available, evict the key in the main thread itself without spawning a goroutine for it
		s.evictKey(key)
	}

	// if slot is acquired, evict the key immediately using goroutine and free up this slot
	go func() {
		s.evictKey(key)
		<-s.evictPool
	}()
}

func (s *store) evictKey(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.keys, key)
}

func parseInput(line string) (string, string, map[string]AttrVal, time.Duration, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", "", nil, 0, fmt.Errorf("invalid input")
	}

	var cmd, key string
	var attrStr []string
	var attrs map[string]AttrVal
	var keyExpiry time.Duration

	cmd = parts[0]
	switch cmd {
	case "get":
		key = parts[1]
		attrStr = nil
	case "put":
		key = parts[1]
		attrStr = parts[3:]
		seconds, _ := strconv.Atoi(parts[2])
		keyExpiry = time.Duration(seconds) * time.Second
	case "delete":
		key = parts[1]
		attrStr = nil
	case "search":
		attrStr = parts[1:]
	case "keys":
		attrStr = nil
	case "exit":
		attrStr = nil
	}

	if attrStr != nil {
		attrs = getAttributes(attrStr)
	}

	return strings.TrimSpace(cmd), strings.TrimSpace(key), attrs, keyExpiry, nil
}

func (s *store) execCmd(cmd string, key string, attrs map[string]AttrVal, keyExpiry time.Duration) {
	switch cmd {
	case "get":
		output, err := s.get(key)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(output.String())
		}
	case "put":
		err := s.put(key, attrs, keyExpiry)
		if err != nil {
			fmt.Println(err)
		}
	case "search":
		for k, v := range attrs {
			fmt.Println(s.search(k, v))
		}
	case "delete":
		s.delete(key)
	case "keys":
		fmt.Println(s.getKeys())
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("invalid command")
	}
}

func getAttributes(attrStr []string) map[string]AttrVal {
	attr := make(map[string]AttrVal)
	for i := 0; i < len(attrStr); i += 2 {
		attr[strings.TrimSpace(attrStr[i])] = parseAttributeValue(strings.TrimSpace(attrStr[i+1]))
	}
	return attr
}

func parseAttributeValue(val string) AttrVal {
	// try int
	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	}

	// try float
	if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
		return floatVal
	}

	// try bool, cannot use strconv.ParseBool(val) because it considers 1 and 0 also as bool
	if val == "true" {
		return true
	} else if val == "false" {
		return false
	}

	// fallback as string
	return val
}

func newValueObject(pairs map[string]AttrVal) *ValueObject {
	types := make(map[string]AttrType)
	for k, v := range pairs {
		types[k] = reflect.TypeOf(v).Kind()
	}
	return &ValueObject{
		types:      types,
		attributes: pairs,
	}
}

func (v *ValueObject) String() string {
	var sb strings.Builder
	first := true
	for k, v := range v.attributes {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(fmt.Sprintf("%s: %v", k, v))
	}
	return sb.String()
}

func String(str []string) string {
	var sb strings.Builder
	first := true
	for i := range str {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(fmt.Sprintf("%s", str[i]))
	}
	return sb.String()
}
