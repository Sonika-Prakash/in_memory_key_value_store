package main

import (
	"fmt"
	"reflect"
	"time"
)

func (s *store) get(key string) (*ValueObject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if output, ok := s.keys[key]; ok {
		return output, nil
	} else {
		return nil, fmt.Errorf("No entry found for %s", key)
	}
}

func (s *store) put(key string, pairs map[string]AttrVal, expiry time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	vo := newValueObject(pairs)

	if currVal, ok := s.keys[key]; ok {
		for k, newVal := range pairs {
			if kind := reflect.TypeOf(newVal).Kind(); kind != currVal.types[k] {
				return fmt.Errorf("Data Type Error for key %s", key)
			}
		}
	}
	s.keys[key] = vo

	if expiry > 0 {
		go func() {
			<-time.After(expiry)
			s.scheduleEviction(key) // schedule this key for eviction after its expiration is reached
		}()
	}

	return nil
}

func (s *store) search(attrKey string, attrVal AttrVal) string {
	res := make([]string, 0)
	for k, valObj := range s.keys {
		if currVal, ok := valObj.attributes[attrKey]; ok && reflect.DeepEqual(currVal, attrVal) {
			res = append(res, k)
		}
	}
	return String(res)
}

func (s *store) delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.keys, key)
}

func (s *store) getKeys() string {
	output := make([]string, 0)
	for k := range s.keys {
		output = append(output, k)
	}
	return String(output)
}
