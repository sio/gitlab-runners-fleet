package tests

import (
	"testing"

	"fmt"
	"strings"

	"scale/gitlab"
)

func TestGraphQL_Echo(t *testing.T) {
	const message string = "endpoint works"
	api := gitlab.API{}
	reply, err := api.GraphQL(
		"query hello($message: String!) {echo(text: $message)}",
		map[string]any{"message": message},
	)
	if err != nil {
		t.Fatal(err)
	}
	value, ok := reply["echo"]
	if !ok {
		t.Fatalf("reply does not contain echo field: %v", reply)
	}
	if !strings.HasSuffix(value.(string), fmt.Sprintf("says: %s", message)) {
		t.Fatalf("unexpected echo output: %q", value)
	}
}

func TestGraphQL_NoParams(t *testing.T) {
	api := gitlab.API{}
	reply, err := api.GraphQL("{currentUser {username}}", nil)
	if err != nil {
		t.Fatal(err)
	}
	name, ok := reply["currentUser"]
	if !ok {
		t.Fatalf("reply does not contain currentUser field: %v", reply)
	}
	if name != nil {
		t.Fatalf("unexpected username in test: %v", reply)
	}
}
