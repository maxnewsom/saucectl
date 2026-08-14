// Package aiauthoring provides domain types for the Sauce Labs AI Test
// Authoring API.
package aiauthoring

// GenerateRequest is the payload for starting a test case generation task.
type GenerateRequest struct {
	Name           string         `json:"name"`
	RunSettings    RunSettings    `json:"runSettings"`
	PromptSettings PromptSettings `json:"promptSettings"`
	Timeout        int            `json:"timeout,omitempty"`
	TestSuiteID    string         `json:"testSuiteId,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
}

// RunSettings describes the target environment a generated test case runs against.
type RunSettings struct {
	Target       Target `json:"target"`
	TestURL      string `json:"testUrl,omitempty"`
	SCTunnelName string `json:"scTunnelName,omitempty"`
}

// Target describes the W3C WebDriver capabilities to run against.
type Target struct {
	Capabilities map[string]interface{} `json:"capabilities"`
}

// PromptSettings describes the natural-language instructions for generation.
type PromptSettings struct {
	Intent   string `json:"intent"`
	MaxSteps int    `json:"maxSteps,omitempty"`
}

// GenerateResponse wraps the result of accepting a generation request.
type GenerateResponse struct {
	Data struct {
		TaskID     string `json:"taskId"`
		SauceJobID string `json:"sauceJobId"`
	} `json:"data"`
}

// ErrorDetail is the standard {code, detail} error shape used throughout the API.
type ErrorDetail struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// GenerationStatus describes the state of an in-flight (or completed) generation task.
type GenerationStatus struct {
	Data struct {
		Status     string       `json:"status"`
		TestCaseID string       `json:"testCaseId,omitempty"`
		Error      *ErrorDetail `json:"error,omitempty"`
	} `json:"data"`
}

// Generation status values.
const (
	StatusQueued     = "QUEUED"
	StatusInProgress = "IN_PROGRESS"
	StatusCompleted  = "COMPLETED"
	StatusFailed     = "FAILED"
)

// TestCase represents a saved AI-authored test case.
type TestCase struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	CreationDate   string   `json:"creationDate,omitempty"`
	LastUpdateDate string   `json:"lastUpdateDate,omitempty"`
	TestSuiteID    string   `json:"testSuiteId,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

// TestCaseList is the paginated response for listing test cases.
type TestCaseList struct {
	Data struct {
		Items []TestCase `json:"items"`
		Total int        `json:"total"`
	} `json:"data"`
}

// TestCaseResponse wraps a single test case.
type TestCaseResponse struct {
	Data TestCase `json:"data"`
}

// RunRequest is the payload for running a saved test case.
type RunRequest struct {
	BuildName    string   `json:"buildName,omitempty"`
	Targets      []Target `json:"targets,omitempty"`
	SCTunnelName string   `json:"scTunnelName,omitempty"`
}

// RunJob describes a single Sauce Labs job created by a test case run (one per target).
type RunJob struct {
	ID         string `json:"id"`
	SauceJobID string `json:"sauceJobId,omitempty"`
	Name       string `json:"name"`
	URL        string `json:"url,omitempty"`
	IsRdc      bool   `json:"isRdc,omitempty"`
	Success    bool   `json:"success,omitempty"`
	Error      string `json:"error,omitempty"`
}

// RunResponse describes a created test case run.
type RunResponse struct {
	Data struct {
		ID           string   `json:"id"`
		TestCaseID   string   `json:"testCaseId"`
		Build        string   `json:"build"`
		Jobs         []RunJob `json:"jobs"`
		CreationDate string   `json:"creationDate"`
	} `json:"data"`
}

// CodeTargetsResponse lists the code-export targets available for a test case.
type CodeTargetsResponse struct {
	Data struct {
		Targets []string `json:"targets"`
	} `json:"data"`
}

// CodeResponse wraps generated source code for a test case.
type CodeResponse struct {
	Data struct {
		Code string `json:"code"`
	} `json:"data"`
}
