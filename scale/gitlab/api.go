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
func (api *API) GraphQL(query string, params map[string]any) (reply graphqlResponse, err error) {
	if api.host == "" {
		api.host = defaultHost
	}
	var url string = fmt.Sprintf("https://%s/api/graphql", api.host)
	for attempts := 0; attempts < apiRetryAttempts; attempts++ {
		reply, err = api.graphql(url, query, params)
		if err == nil {
			return reply, nil
		}
		time.Sleep(apiRetryDelay)
	}
	var zero graphqlResponse
	return zero, err
}

// Execute a straightforward GraphQL API call without any error handling
func (api *API) graphql(url string, query string, params map[string]any) (reply graphqlResponse, err error) {
	var zero graphqlResponse
	var payloadObject = map[string]any{
		"query": query,
	}
	if params != nil {
		payloadObject["variables"] = params
	}
	var payload []byte
	payload, err = json.Marshal(payloadObject)
	if err != nil {
		return zero, fmt.Errorf("could not construct payload: %w", err)
	}

	var body io.Reader = bytes.NewBuffer(payload)
	var req *http.Request
	req, err = http.NewRequest("POST", url, body)
	if err != nil {
		return zero, fmt.Errorf("could not create request object: %w", err)
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
		return zero, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var replyData []byte
	replyData, err = io.ReadAll(resp.Body)
	if err != nil {
		return zero, fmt.Errorf("failed to read HTTP response: %w", err)
	}
	err = json.Unmarshal(replyData, &reply)
	if err != nil {
		return zero, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	if len(reply.Errors) > 0 {
		return reply, fmt.Errorf("GraphQL API returned error(s): %v", reply.Errors)
	}
	if !json.Valid(reply.Data) {
		return zero, fmt.Errorf("reply data is not valid JSON: %s", string(reply.Data))
	}
	err = json.Unmarshal(reply.Data, &reply.Unstructured)
	if err != nil {
		return zero, fmt.Errorf("reply data can not be parsed as JS object: %w", err)
	}
	return reply, nil
}

type graphqlResponse struct {
	// Raw JSON object for custom parsing later
	Data json.RawMessage

	// Straightforward unstructured parsed Data
	Unstructured map[string]any

	// GraphQL error messags
	Errors []map[string]any
}
