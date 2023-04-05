package cloud

import (
	"fmt"
	"net/http"
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
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return nil
}
