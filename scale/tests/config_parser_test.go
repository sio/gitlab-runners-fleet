package tests

import (
	"testing"

	"encoding/json"

	"scale/app"
)

var sampleJson = []byte(`{
	"gitlab_host": "test_host",
	"gitlab_token": "env:TEST_TOKEN",
	"runner_tag": "test_tag",
	"runner_max_jobs": 10,
	"instance_count_max": 11,
	"instance_count_min": 1,
	"instance_grow_max": 3,
	"instance_provision_time": 60,
	"instance_max_age": "12h41m"
}`)

func TestJsonParsing(t *testing.T) {
	var got, want app.Configuration
	var err error
	t.Setenv("TEST_TOKEN", "test_token")
	err = json.Unmarshal(sampleJson, &got)
	if err != nil {
		t.Fatal("failed to parse sampleJson:", err)
	}
	const nanosec = 1000_000_000
	want = app.Configuration{
		GitLabHost:            "test_host",
		GitLabToken:           "test_token",
		RunnerTag:             "test_tag",
		RunnerMaxJobs:         10,
		InstanceCountMax:      11,
		InstanceCountMin:      1,
		InstanceGrowMax:       3,
		InstanceProvisionTime: app.NewDuration(60 * nanosec),
		InstanceMaxAge:        app.NewDuration((12*60*60 + 41*60) * nanosec),
	}
	if got != want {
		t.Fatalf("JSON unmarshaled incorrectly:\n got: %v,\nwant: %v", got, want)
	}
	if got.InstanceProvisionTime.String() != "1m0s" {
		t.Errorf("unexpected string representation for InstanceProvisionTime: %s", got.InstanceProvisionTime.String())
	}
	if got.InstanceMaxAge.String() != "12h41m0s" {
		t.Errorf("unexpected string representation for InstanceMaxAge: %s", got.InstanceMaxAge.String())
	}
}
