package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/saucelabs/saucectl/internal/aiauthoring"
)

// AIAuthoring describes an interface to the AI Test Authoring rest endpoints.
type AIAuthoring struct {
	HTTPClient *retryablehttp.Client
	URL        string
	Username   string
	AccessKey  string
}

// NewAIAuthoring returns a new instance of AIAuthoring. url is expected to be
// the region's API base URL (e.g. https://api.us-west-1.saucelabs.com), the
// ai-authoring path prefix is appended automatically.
func NewAIAuthoring(url string, username string, accessKey string, timeout time.Duration) AIAuthoring {
	return AIAuthoring{
		HTTPClient: NewRetryableClient(timeout),
		URL:        url + "/v1/ai-authoring",
		Username:   username,
		AccessKey:  accessKey,
	}
}

func (c *AIAuthoring) do(ctx context.Context, method, path string, query url.Values, payload interface{}, out interface{}) error {
	reqURL := c.URL + path
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	req, err := NewRetryableRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.SetBasicAuth(c.Username, c.AccessKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusInternalServerError {
		return ErrServerError
	}
	if resp.StatusCode >= http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed; unexpected response code:'%d', msg:'%s'", resp.StatusCode, respBody)
	}
	if out == nil || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// GenerateTestCase kicks off an AI test-case generation task and returns its task ID.
func (c *AIAuthoring) GenerateTestCase(ctx context.Context, req aiauthoring.GenerateRequest) (string, error) {
	var resp aiauthoring.GenerateResponse
	if err := c.do(ctx, http.MethodPost, "/testcases/generate", nil, req, &resp); err != nil {
		return "", err
	}
	return resp.Data.TaskID, nil
}

// GetGenerationStatus returns the current status of a generation task.
func (c *AIAuthoring) GetGenerationStatus(ctx context.Context, taskID string) (aiauthoring.GenerationStatus, error) {
	var resp aiauthoring.GenerationStatus
	err := c.do(ctx, http.MethodGet, "/testcases/generate/"+taskID, nil, nil, &resp)
	return resp, err
}

// ListTestCases returns the test cases visible to the authenticated account.
func (c *AIAuthoring) ListTestCases(ctx context.Context, search string, limit, skip int) ([]aiauthoring.TestCase, error) {
	query := url.Values{}
	if search != "" {
		query.Set("search", search)
	}
	if limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", limit))
	}
	if skip > 0 {
		query.Set("skip", fmt.Sprintf("%d", skip))
	}

	var resp aiauthoring.TestCaseList
	err := c.do(ctx, http.MethodGet, "/testcases", query, nil, &resp)
	return resp.Data.Items, err
}

// GetTestCase returns a single test case by ID.
func (c *AIAuthoring) GetTestCase(ctx context.Context, id string) (aiauthoring.TestCase, error) {
	var resp aiauthoring.TestCaseResponse
	err := c.do(ctx, http.MethodGet, "/testcases/"+id, nil, nil, &resp)
	return resp.Data, err
}

// RunTestCase triggers a run of a saved test case.
func (c *AIAuthoring) RunTestCase(ctx context.Context, id string, req aiauthoring.RunRequest) (aiauthoring.RunResponse, error) {
	var resp aiauthoring.RunResponse
	err := c.do(ctx, http.MethodPost, "/testcases/"+id+"/run", nil, req, &resp)
	return resp, err
}

// GetCodeTargets returns the list of code-export targets valid for a test case.
func (c *AIAuthoring) GetCodeTargets(ctx context.Context, id string) ([]string, error) {
	var resp aiauthoring.CodeTargetsResponse
	err := c.do(ctx, http.MethodGet, "/testcases/"+id+"/code/targets", nil, nil, &resp)
	return resp.Data.Targets, err
}

// GetCode returns the generated source code for a test case, for the given target.
func (c *AIAuthoring) GetCode(ctx context.Context, id, target string) (string, error) {
	query := url.Values{}
	query.Set("target", target)

	var resp aiauthoring.CodeResponse
	err := c.do(ctx, http.MethodGet, "/testcases/"+id+"/code", query, nil, &resp)
	return resp.Data.Code, err
}
