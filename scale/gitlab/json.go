package gitlab

import (
	"fmt"
)

type jsObject = map[string]any

func jsNested(js jsObject, key ...string) (any, error) {
	var current, next jsObject
	var value any
	var ok bool = true
	current = js
	for index, subkey := range key {
		if !ok {
			return nil, fmt.Errorf("type conversion failed for key %s (level %d)", key[index-1], index)
		}
		value, ok = current[subkey]
		if !ok {
			return nil, fmt.Errorf("key not found: %q (level %d)", key[index], index)
		}
		next, ok = value.(jsObject)
		if ok {
			current = next
		}
	}
	return value, nil
}
