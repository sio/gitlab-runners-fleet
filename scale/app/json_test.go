package app

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

// Tests for jsGetAny()
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
	var js map[string]any
	var err error
	err = json.Unmarshal([]byte(input), &js)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		path string
		want any
		fail bool
	}{
		{"hello", "world", false},
		{"foo.bar", "baz", false},
		{"foo.message", map[string]string{"alice": "bob"}, false},
		{"foo.message.alice", "bob", false},
		{"NONEXISTENT", nil, true},
		{"hello.NONEXISTENT", nil, true},
		{"foo.NONEXISTENT", nil, true},
		{"foo.bar.baz.NONEXISTENT", nil, true},
		{"foo.message.alice.NONEXISTENT", nil, true},
	}
	for index, tt := range tests {
		var got any
		got, err = jsGetAny(js, strings.Split(tt.path, ".")...)
		if !tt.fail && err != nil {
			t.Errorf("#%d %v: unexpected error: %v", index, tt.path, err)
			continue
		}
		if tt.fail && err == nil {
			t.Errorf("#%d %v: expected an error, got %v", index, tt.path, got)
			continue
		}
		if got != tt.want && !jsonEqual(got, tt.want) {
			t.Errorf("#%d %v: got %v, want %v", index, tt.path, got, tt.want)
			continue
		}
	}

	stringTests := []struct {
		path string
		want string
		fail bool
	}{
		{"hello", "world", false},
		{"foo.bar", "baz", false},
		{"foo.message.alice", "bob", false},
		{"foo.message", "", true},
		{"NONEXISTENT", "", true},
		{"hello.NONEXISTENT", "", true},
		{"foo.NONEXISTENT", "", true},
		{"foo.bar.baz.NONEXISTENT", "", true},
		{"foo.message.alice.NONEXISTENT", "", true},
	}
	for index, tt := range stringTests {
		var got string
		got, err = JsGet[string](js, strings.Split(tt.path, ".")...)
		if !tt.fail && err != nil {
			t.Errorf("string#%d %v: unexpected error: %v", index, tt.path, err)
			continue
		}
		if tt.fail && err == nil {
			t.Errorf("string#%d %v: expected an error, got %v", index, tt.path, got)
			continue
		}
		if got != tt.want {
			t.Errorf("string#%d %v: got %v, want %v", index, tt.path, got, tt.want)
			continue
		}
	}
}
