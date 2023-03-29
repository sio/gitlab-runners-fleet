package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	userAgent        = "Auto scaling fleet of GitLab runners <https://github.com/sio/gitlab-runners-fleet>"
	defaultHost      = "gitlab.com"
	apiRetryAttempts = 3
	apiRetryDelay    = time.Second * 2
)

var (
	httpClient = http.Client{
		Timeout: time.Second * 3,
	}
)

type jsObject = map[string]any

// Create API object with custom endpoint and/or authorization token.
// Empty strings are accepted for default host / anonymous auth.
func NewAPI(host string, token string) API {
	if host == "" {
		host = defaultHost
	}
	return API{
		host:  host,
		token: token,
	}
}

// GitLab API client.
// Zero value is moderately useful as an anonymous client to gitlab.com
type API struct {
	// GitLab instance hostname
	host string

	// GitLab API token
	token string
}

func (api API) String() string {
	var token string = "none"
	if api.token != "" {
		token = "XXXX"
	}
	return fmt.Sprintf("API{host:%s, token:%s}", api.host, token)
}

// Execute a single GraphQL query with appropriate retries/timeouts
func (api *API) GraphQL(query string, params jsObject) (data jsObject, err error) {
	if api.host == "" {
		api.host = defaultHost
	}
	var url string = fmt.Sprintf("https://%s/api/graphql", api.host)
	for attempts := 0; attempts < apiRetryAttempts; attempts++ {
		data, err = api.graphql(url, query, params)
		if err == nil {
			return data, nil
		}
		time.Sleep(apiRetryDelay)
	}
	return nil, err
}

// Execute a straightforward GraphQL API call without any error handling
func (api *API) graphql(url string, query string, params jsObject) (data jsObject, err error) {
	var payloadObject = jsObject{
		"query": query,
	}
	if params != nil {
		payloadObject["variables"] = params
	}
	var payload []byte
	payload, err = json.Marshal(payloadObject)
	if err != nil {
		return nil, fmt.Errorf("could not construct payload: %w", err)
	}

	fmt.Println(string(payload))

	var body io.Reader = bytes.NewBuffer(payload)
	var req *http.Request
	req, err = http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("could not create request object: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if api.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))
	}

	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var replyData []byte
	replyData, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}
	err = json.Unmarshal(replyData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	var ok, failed bool
	var apiErrors, apiData any
	apiErrors, failed = data["errors"]
	if failed {
		return nil, fmt.Errorf("GraphQL API returned error(s): %v", apiErrors)
	}
	apiData, ok = data["data"]
	if !ok {
		return nil, fmt.Errorf("GraphQL API returned no data: %v", data)
	}
	data, ok = apiData.(jsObject)
	if !ok {
		return nil, fmt.Errorf("type conversion to map[string]any failed: %v", apiData)
	}
	return data, nil
}
