package pbdbapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

// StatusError represents a non-2xx HTTP response from the API.
type StatusError struct {
	StatusCode int
	Body       string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("pbdb status %d: %s", e.StatusCode, strings.TrimSpace(e.Body))
}

// DecodeError wraps JSON decoding failures.
type DecodeError struct {
	Err error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("pbdb decode error: %v", e.Err)
}

func (e *DecodeError) Unwrap() error {
	return e.Err
}

// APIError represents an API-level error encoded in a successful HTTP response.
type APIError struct {
	Message string
}

func (e *APIError) Error() string {
	return "pbdb api error: " + strings.TrimSpace(e.Message)
}

func detectAPIError(body []byte) error {
	if len(body) == 0 {
		return nil
	}

	var probe struct {
		Error  any `json:"error"`
		Errors any `json:"errors"`
	}
	if err := json.Unmarshal(body, &probe); err != nil {
		return nil
	}

	if probe.Error != nil {
		msg := stringifyAny(probe.Error)
		if strings.TrimSpace(msg) != "" {
			return &APIError{Message: msg}
		}
	}

	if probe.Errors != nil {
		if isNonEmptyErrors(probe.Errors) {
			return &APIError{Message: stringifyAny(probe.Errors)}
		}
	}

	return nil
}

func isNonEmptyErrors(v any) bool {
	switch t := v.(type) {
	case []any:
		return len(t) > 0
	case map[string]any:
		return len(t) > 0
	case string:
		return strings.TrimSpace(t) != ""
	default:
		return true
	}
}

func stringifyAny(v any) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}
