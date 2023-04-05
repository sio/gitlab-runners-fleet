package cloud

import (
	"encoding/json"
	"fmt"
	"io"
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

	var raw []byte
	raw, err = io.ReadAll(resp.Body)
	if err != nil {
		return metrics, fmt.Errorf("failed to read HTTP response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return metrics, fmt.Errorf("HTTP %d (%s): %s", resp.StatusCode, req.URL, string(raw))
	}
	err = json.Unmarshal(raw, &metrics)
	if err != nil {
		return metrics, fmt.Errorf("failed to parse JSON metrics for %s: %w\n%s", host.Name, err, string(raw))
	}
	return metrics, nil
}
