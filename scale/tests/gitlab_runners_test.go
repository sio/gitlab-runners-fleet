package tests

import (
	"testing"

	"os"

	"scale/gitlab"
)

func TestRealQuery(t *testing.T) {
	if os.Getenv("GRF_TEST_INTERACTIVE") == "" {
		t.Skip("skipping interactive test")
	}
	var api = gitlab.NewAPI("", os.Getenv("GITLAB_API_TOKEN"))
	err := api.UpdateRunnerAssignments("private-runner")
	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	t.Fatal("This interactive test is set to always fail - to show previous stdout/stderr output")
}
