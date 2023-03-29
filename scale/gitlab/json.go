package gitlab

import (
	"fmt"
	"strings"
)

type jsObject = map[string]any

// Fetch a value from deeply nested JSON
func jsNested(js any, key ...string) (any, error) {
	if len(key) == 0 {
		return js, nil
	}
	var current, next jsObject
	var ok bool
	current, ok = js.(jsObject)
	if !ok {
		return nil, fmt.Errorf("input does not look like a JS object: %v", js)
	}
	var value any
	for index := 0; index < len(key); index++ {
		value, ok = current[key[index]]
		if !ok {
			return nil, fmt.Errorf("key not found: %q (level %d)", strings.Join(key[:index+1], "/"), index)
		}
		next, ok = value.(jsObject)
		if !ok && index < len(key)-1 {
			return nil, fmt.Errorf("type conversion failed for key %s (level %d): %v", strings.Join(key[:index+1], "/"), index, value)
		}
		current = next
	}
	return value, nil
}

// Fetch a string from deeply nested JSON
func jsNestedString(js any, key ...string) (string, error) {
	var value any
	var err error
	value, err = jsNested(js, key...)
	if err != nil {
		return "", err
	}
	var result string
	var ok bool
	result, ok = value.(string)
	if !ok {
		return "", fmt.Errorf("key %v is pointing to non-string value: %v", key, value)
	}
	return result, nil
}
