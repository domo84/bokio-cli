package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type recordedRequest struct {
	method string
	body   string
}

// fakeDoer replays canned responses and records what was actually sent.
type fakeDoer struct {
	responses []*http.Response
	requests  []recordedRequest
}

func (d *fakeDoer) Do(req *http.Request) (*http.Response, error) {
	rec := recordedRequest{method: req.Method}
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		rec.body = string(body)
	}
	d.requests = append(d.requests, rec)

	if len(d.responses) == 0 {
		return nil, errors.New("fakeDoer: no responses left")
	}
	resp := d.responses[0]
	d.responses = d.responses[1:]
	return resp, nil
}

func response(status int, body string, header http.Header) *http.Response {
	if header == nil {
		header = http.Header{}
	}
	return &http.Response{
		StatusCode: status,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func testClient(doer HTTPDoer) *Client {
	return &Client{BaseURL: "https://api.example.test/v1", CompanyID: "c1", HTTP: doer, Token: "token"}
}

// The real payload from a rejected upload. It must surface as a structured error,
// not as raw JSON: the API returns apiError flat, not wrapped in {"error": {...}}.
func TestDoReturnsStructuredAPIError(t *testing.T) {
	const payload = `{"message":"Validation failed with 1 error","code":"validation-error",` +
		`"bokioErrorId":"00000000-0000-4000-8000-000000000001",` +
		`"errors":[{"field":"#/file","message":"Invalid file type. Only application/pdf, image/jpeg and image/png files are allowed."}]}`

	c := testClient(&fakeDoer{responses: []*http.Response{response(400, payload, nil)}})

	var out Upload
	err := c.GetJSON(context.Background(), c.companyURL("/uploads"), &out)
	if err == nil {
		t.Fatal("expected an error")
	}

	var apiErr *BokioError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is %T, want *BokioError: %v", err, err)
	}
	if apiErr.Code != "validation-error" {
		t.Errorf("Code = %q", apiErr.Code)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", apiErr.StatusCode)
	}
	if len(apiErr.Errors) != 1 || apiErr.Errors[0].Field != "#/file" {
		t.Fatalf("Errors = %+v", apiErr.Errors)
	}

	msg := err.Error()
	for _, want := range []string{"validation-error", "#/file", "Invalid file type", "00000000-0000-4000-8000-000000000001"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message %q missing %q", msg, want)
		}
	}
}

func TestParseErrorFallbacks(t *testing.T) {
	cases := []struct {
		name       string
		status     int
		body       string
		wantCode   string
		wantInMsg  string
		wantStatus int
	}{
		{
			name: "oauth error shape", status: 400,
			body:     `{"error":"invalid_grant","error_description":"The code is expired"}`,
			wantCode: "invalid_grant", wantInMsg: "expired", wantStatus: 400,
		},
		{
			name: "non json body", status: 502,
			body: "<html>bad gateway</html>", wantInMsg: "bad gateway", wantStatus: 502,
		},
		{
			name: "empty body falls back to status text", status: 503,
			body: "", wantInMsg: http.StatusText(503), wantStatus: 503,
		},
		{
			name: "message without code", status: 404,
			body: `{"message":"upload not found"}`, wantInMsg: "upload not found", wantStatus: 404,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := parseError(tc.status, []byte(tc.body))

			var apiErr *BokioError
			if !errors.As(err, &apiErr) {
				t.Fatalf("error is %T, want *BokioError", err)
			}
			if apiErr.StatusCode != tc.wantStatus {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tc.wantStatus)
			}
			if tc.wantCode != "" && apiErr.Code != tc.wantCode {
				t.Errorf("Code = %q, want %q", apiErr.Code, tc.wantCode)
			}
			if !strings.Contains(err.Error(), tc.wantInMsg) {
				t.Errorf("message %q missing %q", err.Error(), tc.wantInMsg)
			}
		})
	}
}

// A retried request must resend its body. Replaying a drained body would post
// nothing — silently uploading an empty file or an empty JSON payload.
func TestDoRetriesRateLimitedRequestWithIntactBody(t *testing.T) {
	doer := &fakeDoer{responses: []*http.Response{
		response(429, "", http.Header{"Retry-After": []string{"1"}}),
		response(200, `{"id":"upload-1","description":"receipt"}`, nil),
	}}
	c := testClient(doer)

	var out Upload
	err := c.PostJSON(context.Background(), c.companyURL("/uploads"),
		map[string]string{"description": "receipt"}, &out)
	if err != nil {
		t.Fatalf("PostJSON: %v", err)
	}
	if out.ID != "upload-1" {
		t.Errorf("ID = %q, want upload-1", out.ID)
	}

	if len(doer.requests) != 2 {
		t.Fatalf("made %d requests, want 2", len(doer.requests))
	}
	if doer.requests[1].body == "" {
		t.Fatal("retry sent an empty body")
	}
	if doer.requests[0].body != doer.requests[1].body {
		t.Errorf("body not replayed: first %q, retry %q",
			doer.requests[0].body, doer.requests[1].body)
	}
}

func TestDoReturnsRateLimitErrorWhenRetriesExhausted(t *testing.T) {
	original := maxRateLimitRetries
	maxRateLimitRetries = 0 // no retries, so the test does not sleep
	defer func() { maxRateLimitRetries = original }()

	c := testClient(&fakeDoer{responses: []*http.Response{
		response(429, `{"code":"rate-limit-exceeded","message":"Too many requests"}`, nil),
	}})

	var out Upload
	err := c.GetJSON(context.Background(), c.companyURL("/uploads"), &out)

	var apiErr *BokioError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error is %T, want *BokioError: %v", err, err)
	}
	if !apiErr.IsRateLimited() {
		t.Errorf("StatusCode = %d, want 429", apiErr.StatusCode)
	}
}

func TestDoAbortsWaitWhenContextCancelled(t *testing.T) {
	doer := &fakeDoer{responses: []*http.Response{
		response(429, "", http.Header{"Retry-After": []string{"60"}}),
	}}
	c := testClient(doer)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var out Upload
	err := c.GetJSON(ctx, c.companyURL("/uploads"), &out)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if len(doer.requests) != 1 {
		t.Errorf("made %d requests, want 1", len(doer.requests))
	}
}

func TestRetryAfter(t *testing.T) {
	cases := []struct {
		name   string
		header http.Header
		want   time.Duration
	}{
		{
			name:   "bokio header wins over standard",
			header: http.Header{"Bokio-Ratelimit-Retryafter": []string{"7"}, "Retry-After": []string{"30"}},
			want:   7 * time.Second,
		},
		{
			name:   "falls back to standard header",
			header: http.Header{"Retry-After": []string{"12"}},
			want:   12 * time.Second,
		},
		{
			name:   "clamps absurd values",
			header: http.Header{"Retry-After": []string{"86400"}},
			want:   maxRetryWaitSeconds * time.Second,
		},
		{
			name:   "ignores unparseable values",
			header: http.Header{"Retry-After": []string{"Wed, 21 Oct 2026 07:28:00 GMT"}},
			want:   defaultRetryWait,
		},
		{name: "defaults when absent", header: http.Header{}, want: defaultRetryWait},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := retryAfter(tc.header); got != tc.want {
				t.Errorf("retryAfter() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRewind(t *testing.T) {
	t.Run("replays an in-memory body", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
			"https://api.example.test/v1/x", strings.NewReader(`{"a":1}`))
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		if _, err := io.ReadAll(req.Body); err != nil { // drain, as a send would
			t.Fatalf("drain: %v", err)
		}

		retryReq, err := rewind(req)
		if err != nil {
			t.Fatalf("rewind: %v", err)
		}
		body, err := io.ReadAll(retryReq.Body)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		if string(body) != `{"a":1}` {
			t.Errorf("body = %q", body)
		}
	})

	t.Run("passes through a bodyless request", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet,
			"https://api.example.test/v1/x", nil)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		if _, err := rewind(req); err != nil {
			t.Fatalf("rewind: %v", err)
		}
	})

	t.Run("errors when the body cannot be replayed", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
			"https://api.example.test/v1/x", io.NopCloser(strings.NewReader("stream")))
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}
		req.GetBody = nil // an opaque stream, as http would leave it
		if _, err := rewind(req); err == nil {
			t.Error("expected an error for a non-replayable body")
		}
	})
}
