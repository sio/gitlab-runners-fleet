//
// Manage project-level runners (the only runner type available to free users at gitlab.com)
//

package gitlab

import (
	"encoding/json"
	"fmt"
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
		runners  = make(map[string]string)
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
				runners[runner.Node.ID] = runner.Node.Status
			}
		}
		if !reply.Projects.PageInfo.HasNextPage {
			break
		}
		page = reply.Projects.PageInfo.EndCursor
	}
	fmt.Println(string(graphql.Data))
	fmt.Println("Runners:", runners)
	fmt.Println("Projects:", projects)
	return nil
}

//
//func removeRunner(uid string) {
//}
//
//func assignRunner(runnerUID string, projectUID []string) {
//}

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
					Node struct {
						ID     string
						Status string
					}
				}
			}
		}
	}
}

type void struct{}
