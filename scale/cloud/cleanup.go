package cloud

import (
	"fmt"
	"net/http"
	"time"
)

func (fleet *Fleet) Cleanup(host *Host) (err error) {
	defer func() { host.Status = Destroying }()
	if fleet.Entrypoint == "" {
		return fmt.Errorf("HTTP endpoint is not defined")
	}
	var (
		req  *http.Request
		resp *http.Response
	)
	req, err = http.NewRequest("POST", fmt.Sprintf("http://%s/unregister", fleet.Entrypoint), nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Host = host.Name
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, req.URL)
	}
	return nil
}
