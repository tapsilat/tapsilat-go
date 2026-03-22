package tapsilat

import (
	"encoding/json"
	"fmt"
)

type APIError struct {
	StatusCode int
	Status     string
	Code       string
	Message    string
	RawBody    string
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
	if json.Unmarshal(body, &payload) != nil {
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
	} else if value, ok := payload["error"]; ok {
		err.Message = fmt.Sprint(value)
	}

	return err
}
