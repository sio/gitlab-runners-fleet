package tests

import (
	"testing"

	"encoding/json"
	"os"
	"path/filepath"

	"scale/cloud"
)

func TestSerialization(t *testing.T) {
	var original, restored *cloud.Fleet
	original = &cloud.Fleet{}
	var err error

	var data = []byte(`
	{
		"entrypoint": "hello.world.tld",
		"hosts": [
			{"name": "first-host", "status": 4},
			{"name": "second-host", "status": 2, "created_at": "2022-09-29T12:23:34Z"}
		]
	}
	`)
	if err = json.Unmarshal(data, original); err != nil {
		t.Fatal(err)
	}

	temp, err := os.MkdirTemp("", "scaler-state")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(temp)

	stateFile := filepath.Join(temp, "state.json")
	original.Save(stateFile)

	var saved []byte
	saved, err = os.ReadFile(stateFile)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(saved))

	restored = &cloud.Fleet{}
	restored.Load(stateFile)
	if original.Entrypoint != restored.Entrypoint {
		t.Errorf("entrypoint was %q, became %q", original.Entrypoint, restored.Entrypoint)
	}
	if len(restored.Hosts) != len(original.Hosts) {
		t.Fatalf("saved %d hosts, restored %d hosts", len(original.Hosts), len(restored.Hosts))
	}
	for key := range original.Hosts {
		var got, want cloud.Host
		want = *original.Hosts[key]
		got = *restored.Hosts[key]
		if got != want {
			t.Errorf("serialization failure for host %q: got %v, want %v", key, got, want)
		}
	}
	if len(original.Hosts) == 0 {
		t.Errorf("empty original host list: %v", original)
	}
	if len(restored.Hosts) == 0 {
		t.Errorf("empty restored host list: %v", restored)
	}
}
