package api

import "fmt"

// BokioError represents an error response from the Bokio API.
type BokioError struct {
	Code         string            `json:"code"`
	InnerCode    string            `json:"innerCode,omitempty"`
	Message      string            `json:"message"`
	BokioErrorID string            `json:"bokioErrorId,omitempty"`
	Errors       []ValidationError `json:"errors,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error BokioError `json:"error"`
}

func (e *BokioError) Error() string {
	msg := fmt.Sprintf("[%s] %s", e.Code, e.Message)
	if e.BokioErrorID != "" {
		msg += fmt.Sprintf(" (bokioErrorId: %s)", e.BokioErrorID)
	}
	for _, v := range e.Errors {
		msg += fmt.Sprintf("\n  %s: %s", v.Field, v.Message)
	}
	return msg
}
