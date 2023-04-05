package cloud

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Metrics struct {
	JobsTotal   int `json:"gitlab_runner_jobs_total"`
	JobsRunning int `json:"gitlab_runner_jobs"`
}

func (fleet *Fleet) Metrics(host *Host) (Metrics, error) {
	var (
		err     error
		metrics Metrics
		req     *http.Request
		resp    *http.Response
	)
	if fleet.Entrypoint == "" {
		return metrics, fmt.Errorf("HTTP endpoint is not yet known")
	}
	req, err = http.NewRequest("GET", fmt.Sprintf("http://%s/metrics", fleet.Entrypoint), nil)
	if err != nil {
		return metrics, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Host = host.Name
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return metrics, fmt.Errorf("metrics request failed for %s: %w", host.Name, err)
	}
	defer func() { _ = resp.Body.Close() }()

	err = json.NewDecoder(resp.Body).Decode(&metrics)
	if err != nil {
		return metrics, fmt.Errorf("failed to parse host metrics for %s: %w", host.Name, err)
	}
	return metrics, nil
}
