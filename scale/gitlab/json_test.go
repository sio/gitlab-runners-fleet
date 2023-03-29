package gitlab

import (
	"testing"

	"bytes"
	"encoding/json"
	"strings"
)

// Check if two inputs produce the same JSON
func jsonEqual(a, b any) bool {
	var jsonA, jsonB []byte
	var errA, errB error
	jsonA, errA = json.Marshal(a)
	jsonB, errB = json.Marshal(b)
	if errA == nil && errB == nil {
		return bytes.Equal(jsonA, jsonB)
	}
	return false
}

// Tests for jsNested()
func TestNestedAccess(t *testing.T) {
	input := `{
		"hello": "world",
		"foo": {
			"bar": "baz",
			"message": {
				"alice": "bob"
			}
		}
	}`
	var js jsObject
	var err error
	err = json.Unmarshal([]byte(input), &js)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		path  string
		want  any
		iserr bool
	}{
		{"hello", "world", false},
		{"foo:bar", "baz", false},
		{"foo:message", map[string]string{"alice": "bob"}, false},
		{"foo:message:alice", "bob", false},
		{"NONEXISTENT", nil, true},
		{"hello:NONEXISTENT", nil, true},
		{"foo:NONEXISTENT", nil, true},
		{"foo:bar:baz:NONEXISTENT", nil, true},
		{"foo:message:alice:NONEXISTENT", nil, true},
	}
	for index, tt := range tests {
		var got any
		got, err = jsNested(js, strings.Split(tt.path, ":")...)
		if !tt.iserr && err != nil {
			t.Errorf("#%d %v: unexpected error: %v", index, tt.path, err)
		}
		if tt.iserr && err == nil {
			t.Errorf("#%d %v: expected an error, got %v", index, tt.path, got)
		}
		if got != tt.want && !jsonEqual(got, tt.want) {
			t.Errorf("#%d %v: got %v, want %v", index, tt.path, got, tt.want)
		}
	}
}
