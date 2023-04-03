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

const second = 1_000_000_000 // nanoseconds

func TestJsonParsing(t *testing.T) {
	var got, want app.Configuration
	var err error
	t.Setenv("TEST_TOKEN", "test_token")
	err = json.Unmarshal(sampleJson, &got)
	if err != nil {
		t.Fatal("failed to parse sampleJson:", err)
	}
	want = app.Configuration{
		GitLabHost:            "test_host",
		GitLabToken:           "test_token",
		RunnerTag:             "test_tag",
		RunnerMaxJobs:         10,
		InstanceCountMax:      11,
		InstanceCountMin:      1,
		InstanceGrowMax:       3,
		InstanceProvisionTime: app.NewDuration(60 * second),
		InstanceMaxAge:        app.NewDuration((12*60*60 + 41*60) * second),
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

func TestDefaultConfig(t *testing.T) {
	var err error
	var got, want app.Configuration

	got = app.DefaultConfiguration
	err = json.Unmarshal([]byte(`
		{"gitlab_host": "test_host", "gitlab_token": "test_token"}
	`), &got)
	if err != nil {
		t.Fatal("failed to parse JSON input:", err)
	}

	want = app.Configuration{
		ScalerState:           "scaler.state",
		TerraformState:        "terraform.tfstate",
		GitLabHost:            "test_host",
		GitLabToken:           "test_token",
		RunnerTag:             "gitlab_runner_fleet",
		RunnerMaxJobs:         3,
		InstanceCountMax:      10,
		InstanceCountMin:      0,
		InstanceGrowMax:       3,
		InstanceProvisionTime: app.NewDuration(10 * 60 * second),
		InstanceMaxAge:        app.NewDuration(24 * 60 * 60 * second),
		InstanceMaxIdleTime:   app.NewDuration(40 * 60 * second),
	}
	if got != want {
		t.Fatalf("JSON unmarshaled incorrectly:\n got: %v,\nwant: %v", got, want)
	}
	if got == app.DefaultConfiguration {
		t.Fatalf("DeafaulConfiguration got modified during unmarshaling!")
	}
}
