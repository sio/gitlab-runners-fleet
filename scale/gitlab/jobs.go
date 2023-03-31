package gitlab

import (
	"encoding/json"
	"log"
	"strings"
)

func (api *API) CountPendingJobs() int {
	var page string
	var count int
	for {
		var gql graphqlResponse
		var err error
		gql, err = api.GraphQL(queryPendingJobs, map[string]any{"page": page})
		if err != nil {
			log.Printf("failed to list pending jobs (page %q): %v", page, err)
			return count
		}
		var reply jobsResult
		err = json.Unmarshal(gql.Data, &reply)
		if err != nil {
			log.Printf("failed to parse API response: %v\n%s", err, string(gql.Data))
			return count
		}
		count += reply.Count("pending")
		page = reply.Projects.PageInfo.EndCursor
		if !reply.Projects.PageInfo.HasNextPage {
			break
		}
	}
	return count
}

type jobsResult struct {
	Projects struct {
		PageInfo struct {
			HasNextPage bool
			EndCursor   string
		}
		Nodes []struct {
			Running pipelineList
			Waiting pipelineList
			Pending pipelineList
		}
	}
}

func (result jobsResult) Count(status string) int {
	var count int
	for _, project := range result.Projects.Nodes {
		count += project.Running.Count(status)
		count += project.Waiting.Count(status)
		count += project.Pending.Count(status)
	}
	return count
}

type pipelineList struct {
	Nodes []struct {
		Jobs struct {
			Nodes []struct {
				Status string
			}
		}
	}
}

func (list pipelineList) Count(status string) int {
	status = strings.ToLower(status)
	var count int
	for _, pipeline := range list.Nodes {
		for _, job := range pipeline.Jobs.Nodes {
			if strings.ToLower(job.Status) == status {
				count++
			}
		}
	}
	return count
}
