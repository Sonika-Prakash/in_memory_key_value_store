package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func newStore() *store {
	return &store{
		keys: make(map[string]*ValueObject),
	}
}

func parseInput(line string) (string, string, map[string]AttrVal, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "", "", nil, fmt.Errorf("invalid input")
	}

	var cmd, key string
	var attrStr []string
	var attrs map[string]AttrVal

	cmd = parts[0]
	switch cmd {
	case "get":
		key = parts[1]
		attrStr = nil
	case "put":
		key = parts[1]
		attrStr = parts[2:]
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

	return strings.TrimSpace(cmd), strings.TrimSpace(key), attrs, nil
}

func (s *store) execCmd(cmd string, key string, attrs map[string]AttrVal) {
	switch cmd {
	case "get":
		output, err := s.get(key)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(output.String())
		}
	case "put":
		err := s.put(key, attrs)
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
