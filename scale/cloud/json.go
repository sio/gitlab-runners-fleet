package cloud

import (
	"fmt"
	"strings"
)

// Fetch vaule of a specific type from deeply nested JSON
func JsGet[T any](tree any, key ...string) (T, error) {
	var value any
	var err error
	value, err = jsGetAny(tree, key...)
	var zero T
	if err != nil {
		return zero, err
	}
	var result T
	var ok bool
	result, ok = value.(T)
	if !ok {
		return zero, fmt.Errorf("type error: value for %s is not %T but %T: %v", repr(key), zero, value, value)
	}
	return result, nil
}

// Fetch a value from deeply nested JSON
func jsGetAny(tree any, key ...string) (any, error) {
	if len(key) == 0 {
		return tree, nil
	}
	var current, next map[string]any
	var ok bool
	current, ok = tree.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("input does not look like a parsed JSON: %v", tree)
	}
	var value any
	for index := 0; index < len(key); index++ {
		value, ok = current[key[index]]
		if !ok {
			return nil, fmt.Errorf("key not found: %s (level %d)", repr(key[:index+1]), index)
		}
		next, ok = value.(map[string]any)
		if !ok && index < len(key)-1 {
			return nil, fmt.Errorf("can not go deeper than key %s (level %d): %v", repr(key[:index+1]), index, value)
		}
		current = next
	}
	return value, nil
}

// String representation of deep JSON key for error messages
func repr(key []string) string {
	const pathSeparator = "."
	return strings.Join(key, pathSeparator)
}
