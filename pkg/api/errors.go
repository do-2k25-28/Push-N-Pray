package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPError represents an API error response.
type HTTPError struct {
	StatusCode int
	Status     string
	Message    string
}

func (e *HTTPError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("api: http %d: %s", e.StatusCode, e.Status)
	}
	return fmt.Sprintf("api: http %d: %s: %s", e.StatusCode, e.Status, e.Message)
}

type errorResponse struct {
	Message string `json:"message"`
}

func readHTTPError(resp *http.Response) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &HTTPError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	if len(data) == 0 {
		return &HTTPError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	var payload errorResponse
	if err := json.Unmarshal(data, &payload); err == nil && payload.Message != "" {
		return &HTTPError{StatusCode: resp.StatusCode, Status: resp.Status, Message: payload.Message}
	}

	return &HTTPError{StatusCode: resp.StatusCode, Status: resp.Status, Message: string(data)}
}
