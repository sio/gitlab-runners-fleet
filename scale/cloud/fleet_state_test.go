package cloud

import (
	"testing"

	"encoding/json"
	"os"
	"path/filepath"
)

func TestSerialization(t *testing.T) {
	var original, restored *Fleet
	original = &Fleet{}
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

	restored = &Fleet{}
	restored.Load(stateFile)
	if original.entrypoint != restored.entrypoint {
		t.Errorf("entrypoint was %q, became %q", original.entrypoint, restored.entrypoint)
	}
	if len(restored.hosts) != len(original.hosts) {
		t.Fatalf("saved %d hosts, restored %d hosts", len(original.hosts), len(restored.hosts))
	}
	for key := range original.hosts {
		var got, want Host
		want = *original.hosts[key]
		got = *restored.hosts[key]
		if got != want {
			t.Errorf("serialization failure for host %q: got %v, want %v", key, got, want)
		}
	}
	if len(original.hosts) == 0 {
		t.Errorf("empty original host list: %v", original)
	}
	if len(restored.hosts) == 0 {
		t.Errorf("empty restored host list: %v", restored)
	}
}
