package tests

import (
	"testing"

	"os"
	"path/filepath"

	"scale/cloud"
)

func TestSerialization(t *testing.T) {
	original := cloud.Fleet{
		Entrypoint: "hello",
		Hosts: []cloud.Host{
			{Name: "foo"},
			{Name: "bar", Status: cloud.Provisioning},
		},
	}
	temp, err := os.MkdirTemp("", "scaler-state")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(temp)

	stateFile := filepath.Join(temp, "state.json")
	original.Save(stateFile)
	var restored cloud.Fleet
	restored.Load(stateFile)
	if original.Entrypoint != restored.Entrypoint {
		t.Errorf("entrypoint was %s, became %s", original.Entrypoint, restored.Entrypoint)
	}
	if len(restored.Hosts) != len(original.Hosts) {
		t.Fatalf("saved %d hosts, restored %d hosts", len(original.Hosts), len(restored.Hosts))
	}
	for i := 0; i < len(original.Hosts); i++ {
		var got, want cloud.Host
		want = original.Hosts[i]
		got = restored.Hosts[i]
		if got != want {
			t.Errorf("serialization failure: got %v, want %v", got, want)
		}
	}
}
