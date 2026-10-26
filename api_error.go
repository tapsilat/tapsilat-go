package tapsilat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type APIError struct {
	StatusCode int
	Status     string
	Code       string
	Message    string
	// Errors holds per-field validation messages when the server returns them.
	Errors map[string]any
	// Retryable is set when the error object says the request can be retried, as Elements errors do.
	Retryable bool
	RawBody   string
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Code != "" && e.Message != "" {
		return fmt.Sprintf("API request failed with status %d (%s): [%s] %s", e.StatusCode, e.Status, e.Code, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("API request failed with status %d (%s): %s", e.StatusCode, e.Status, e.Message)
	}
	if e.RawBody != "" {
		return fmt.Sprintf("API request failed with status %d (%s): %s", e.StatusCode, e.Status, e.RawBody)
	}
	return fmt.Sprintf("API request failed with status %d (%s)", e.StatusCode, e.Status)
}

func newAPIError(statusCode int, status string, body []byte) *APIError {
	err := &APIError{
		StatusCode: statusCode,
		Status:     status,
		RawBody:    string(body),
	}

	var payload map[string]any
	decode := json.NewDecoder(bytes.NewReader(body))
	decode.UseNumber()
	if decode.Decode(&payload) != nil {
		return err
	}

	if value, ok := payload["status"]; ok {
		err.Status = fmt.Sprint(value)
	}
	if value, ok := payload["code"]; ok {
		err.Code = fmt.Sprint(value)
	}
	if value, ok := payload["message"]; ok {
		err.Message = fmt.Sprint(value)
	} else if nested, ok := payload["error"].(map[string]any); ok {
		// Elements routes nest the error: {"error": {"code", "message", "retryable"}}.
		if value, ok := nested["code"]; ok && err.Code == "" {
			err.Code = fmt.Sprint(value)
		}
		if value, ok := nested["message"]; ok {
			err.Message = fmt.Sprint(value)
		}
		err.Retryable, _ = nested["retryable"].(bool)
	} else if value, ok := payload["error"]; ok {
		err.Message = fmt.Sprint(value)
	}
	if value, ok := payload["errors"].(map[string]any); ok {
		err.Errors = value
	}

	return err
}

// isHTMLResponse reports whether a successful response carries an HTML page instead of JSON.
func isHTMLResponse(contentType string, body []byte) bool {
	return strings.Contains(strings.ToLower(contentType), "text/html") || bytes.HasPrefix(body, []byte("<"))
}

func newUnexpectedResponseError(statusCode int, status, contentType string, body []byte) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Status:     status,
		Message:    fmt.Sprintf("unexpected non-JSON response (Content-Type %q); check that the endpoint URL and route are correct", contentType),
		RawBody:    string(body),
	}
}
