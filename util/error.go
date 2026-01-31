package util

import "fmt"

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

// helper to create
func NewAPIError(status int, msg string) *APIError {
	return &APIError{
		StatusCode: status,
		Message:    msg,
	}
}
