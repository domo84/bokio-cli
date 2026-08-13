package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// BokioError represents an error response from the Bokio API.
//
// The API returns this payload flat at the top level — {"code": ..., "message": ...}
// — not wrapped in an {"error": {...}} object.
type BokioError struct {
	// StatusCode is the HTTP status. It is not part of the payload.
	StatusCode   int               `json:"-"`
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

func (e *BokioError) Error() string {
	var b strings.Builder

	switch {
	case e.Code != "":
		fmt.Fprintf(&b, "[%s] ", e.Code)
	case e.StatusCode != 0:
		fmt.Fprintf(&b, "[HTTP %d] ", e.StatusCode)
	}
	b.WriteString(e.Message)

	if e.InnerCode != "" {
		fmt.Fprintf(&b, " (innerCode: %s)", e.InnerCode)
	}
	if e.BokioErrorID != "" {
		fmt.Fprintf(&b, " (bokioErrorId: %s)", e.BokioErrorID)
	}
	for _, v := range e.Errors {
		fmt.Fprintf(&b, "\n  %s: %s", v.Field, v.Message)
	}
	return b.String()
}

// IsNotFound reports whether the API answered 404.
func (e *BokioError) IsNotFound() bool { return e.StatusCode == http.StatusNotFound }

// IsRateLimited reports whether the API answered 429.
func (e *BokioError) IsRateLimited() bool { return e.StatusCode == http.StatusTooManyRequests }

// parseError turns a failed response into an error, preserving the API's error code
// and per-field validation messages. Bodies that are not a recognisable error payload
// fall back to the raw text so nothing is lost.
func parseError(statusCode int, body []byte) error {
	var apiErr BokioError
	if err := json.Unmarshal(body, &apiErr); err == nil && (apiErr.Code != "" || apiErr.Message != "") {
		apiErr.StatusCode = statusCode
		return &apiErr
	}

	// The OAuth endpoints use {"error": ..., "error_description": ...} instead.
	var oauthErr struct {
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &oauthErr); err == nil && oauthErr.Error != "" {
		return &BokioError{
			StatusCode: statusCode,
			Code:       oauthErr.Error,
			Message:    oauthErr.Description,
		}
	}

	message := strings.TrimSpace(string(body))
	if message == "" {
		message = http.StatusText(statusCode)
	}
	return &BokioError{StatusCode: statusCode, Message: message}
}
