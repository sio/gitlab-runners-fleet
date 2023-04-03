package tests

import (
	"testing"

	"time"

	"scale/cloud"
)

func TestReadTerraformState(t *testing.T) {
	const filename = "fleet_tfstate.sample"
	var fleet cloud.Fleet
	err := fleet.LoadTerraformState(filename)
	if err != nil {
		t.Fatal(err)
	}
	var size = len(fleet.Hosts())
	if size != 2 {
		t.Errorf("incorrect number of host entries: got %d, want %d\n%v", size, 2, fleet.Hosts())
	}
	var zero time.Time
	for _, name := range []string{"hello", "world"} {
		host, ok := fleet.Get(name)
		if !ok {
			t.Errorf("host not found: %s", name)
			continue
		}

		if host.CreatedAt == zero {
			t.Errorf("host.CreatedAt field is empty: %s", host)
		}
	}
}
