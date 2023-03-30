package tests

import (
	"testing"

	"os"

	"scale/gitlab"
)

func TestRealQuery(t *testing.T) {
	var api = gitlab.NewAPI("", os.Getenv("GITLAB_API_TOKEN"))
	err := api.UpdateRunnerAssignments("private-runner")
	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	t.Fatal("OK")
}
