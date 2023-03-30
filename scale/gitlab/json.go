package gitlab

import (
	"fmt"
	"strings"
)

// Fetch a value from deeply nested JSON
func jsGetAny(js any, key ...string) (any, error) {
	if len(key) == 0 {
		return js, nil
	}
	var current, next map[string]any
	var ok bool
	current, ok = js.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("input does not look like a JS object: %v", js)
	}
	var value any
	for index := 0; index < len(key); index++ {
		value, ok = current[key[index]]
		if !ok {
			return nil, fmt.Errorf("key not found: %q (level %d)", strings.Join(key[:index+1], "/"), index)
		}
		next, ok = value.(map[string]any)
		if !ok && index < len(key)-1 {
			return nil, fmt.Errorf("type conversion failed for key %s (level %d): %v", strings.Join(key[:index+1], "/"), index, value)
		}
		current = next
	}
	return value, nil
}

// Fetch a string from deeply nested JSON
func jsGetString(js any, key ...string) (string, error) {
	return jsGet[string](js, key...)
}

// Fetch a specific type from deeply nested JSON
func jsGet[T any](js any, key ...string) (T, error) {
	var value any
	var err error
	value, err = jsGetAny(js, key...)
	var zero T
	if err != nil {
		return zero, err
	}
	var result T
	var ok bool
	result, ok = value.(T)
	if !ok {
		return zero, fmt.Errorf("type error: value for %v is not %T but %T = %v", key, zero, value, value)
	}
	return result, nil
}
