//
// Manage project-level runners (the only runner type available to free users at gitlab.com)
//

package gitlab

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Execute all maintenance chores:
//   - Purge dead runner entries
//   - Assign new runners to all projects
func (api API) UpdateRunnerAssignments(tag string) error {
	var (
		err      error
		graphql  graphqlResponse
		reply    projectQueryReply
		page     string
		projects = make(map[string]void)
		runners  = make(map[string]*runnerInfo)
	)

	// Find all projects that current user is a member of
	// and all available project runners by tag
	for {
		graphql, err = api.GraphQL(queryRunnersByProject, map[string]any{"page": page, "tag": tag})
		if err != nil {
			return fmt.Errorf("projects query failed (page %q): %w", page, err)
		}
		err = json.Unmarshal(graphql.Data, &reply)
		if err != nil {
			return fmt.Errorf("failed to parse GraphQL reply: %w", err)
		}
		for _, project := range reply.Projects.Nodes {
			projects[project.ID] = void{}
			for _, runner := range project.Runners.Edges {
				r, ok := runners[runner.Node.ID]
				if !ok {
					var newRunner runnerInfo = runner.Node
					runners[newRunner.ID] = &newRunner
					r = &newRunner
				}
				// ProjectCount is intentionally calculated manually and not
				// taken from API.
				// Because of runner/project visibility settings not all
				// projects that the runner is assigned to might be visible to
				// current user.
				// We want to check how many visible projects are serviced
				// by this runner, not how many total projects are.
				r.IncrementCount()
			}
		}
		if !reply.Projects.PageInfo.HasNextPage {
			break
		}
		page = reply.Projects.PageInfo.EndCursor
	}

	// Assign detected runners to all projects
	var projectGID = make([]string, 0, len(projects))
	for gid := range projects {
		projectGID = append(projectGID, gid)
	}
	for _, runner := range runners {
		if strings.ToLower(runner.Status) != "online" {
			log.Printf("unregistering runner %s", runner)
			err = api.unregisterRunner(runner.ID)
			if err != nil {
				log.Printf("could not unregister runner %s: %v", runner, err)
			}
			continue
		}
		if runner.ProjectCount == len(projectGID) {
			continue
		}
		log.Printf("updating project assignments for %s", runner)
		err = api.assignRunner(runner.ID, projectGID)
		if err != nil {
			log.Printf("could not reassign runner %s: %v", runner, err)
		}
		// TODO: wait for a while after reassigning runners to allow them to pick up jobs
	}
	return nil
}

// Assign runner to multiple projects at once
func (api *API) assignRunner(runnerGID string, projectGID []string) error {
	var graphql graphqlResponse
	var err error
	graphql, err = api.GraphQL(mutationRunnerAssignment, map[string]any{"runner": runnerGID, "projects": projectGID})
	if err != nil {
		return fmt.Errorf("runner assignment failed: %w", err)
	}
	var reply runnerAssignmentReply
	err = json.Unmarshal(graphql.Data, &reply)
	if err != nil {
		return fmt.Errorf("runner assignment completed with unparsable reply:%w\n%s", err, string(graphql.Data))
	}
	if len(reply.RunnerUpdate.Errors) > 0 {
		return fmt.Errorf("runner assignment completed with errors: %v", reply.RunnerUpdate.Errors)
	}
	return nil
}

// Remove runner from GitLab API
func (api *API) unregisterRunner(gid string) (err error) {
	err = api.assignRunner(gid, []string{})
	if err != nil {
		log.Printf("failed to clear project assignments for %s: %v", gid, err)
	}

	var graphql graphqlResponse
	graphql, err = api.GraphQL(mutationRunnerRemove, map[string]any{"runner": gid})
	if err != nil {
		return fmt.Errorf("failed to unregister %s: %w", gid, err)
	}
	var reply runnerRemovalReply
	err = json.Unmarshal(graphql.Data, &reply)
	if err != nil {
		return fmt.Errorf("runner removal completed with unparsable reply:%w\n%s", err, string(graphql.Data))
	}
	if len(reply.RunnerDelete.Errors) > 0 {
		return fmt.Errorf("runner removal completed with errors: %v", reply.RunnerDelete.Errors)
	}
	return nil
}

type runnerInfo struct {
	ID           string
	Status       string
	Description  string
	ProjectCount int // intentionally not queried via API
}

func (r *runnerInfo) IncrementCount() {
	r.ProjectCount++
}

func (r *runnerInfo) String() string {
	return fmt.Sprintf("%s;%s;%s", r.Description, r.ID, r.Status)
}

type projectQueryReply struct {
	Projects struct {
		PageInfo struct {
			HasNextPage bool
			EndCursor   string
		}
		Nodes []struct {
			ID      string
			Runners struct {
				Edges []struct {
					Node runnerInfo
				}
			}
		}
	}
}

type runnerAssignmentReply struct {
	RunnerUpdate struct {
		Errors []map[string]any
	}
}

type runnerRemovalReply struct {
	RunnerDelete struct {
		Errors []map[string]any
	}
}

type void struct{}
