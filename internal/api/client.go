package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://api.bokio.se/v1"

	// defaultRetryWait is used when a 429 carries no usable wait hint.
	defaultRetryWait = 5 * time.Second
	// maxRetryWaitSeconds caps the honoured wait so an out-of-range header value
	// cannot hang the command.
	maxRetryWaitSeconds = 60
)

// maxRateLimitRetries bounds how many times a 429 is retried before the error is
// returned to the caller. A variable so tests can exercise the exhaustion path.
var maxRateLimitRetries = 3

// HTTPDoer abstracts http.Client for testing.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the Bokio API client.
type Client struct {
	BaseURL   string
	CompanyID string
	HTTP      HTTPDoer
	Token     string
}

// NewClient creates a new API client.
func NewClient(token, companyID string) *Client {
	return &Client{
		BaseURL:   DefaultBaseURL,
		CompanyID: companyID,
		HTTP:      &http.Client{Timeout: 30 * time.Second},
		Token:     token,
	}
}

// companyURL builds a URL under /companies/{companyId}.
func (c *Client) companyURL(path string) string {
	return fmt.Sprintf("%s/companies/%s%s", c.BaseURL, c.CompanyID, path)
}

// generalURL builds a URL at the API root.
func (c *Client) generalURL(path string) string {
	return fmt.Sprintf("%s%s", c.BaseURL, path)
}

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	if body != nil && method != http.MethodGet {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	for attempt := 0; ; attempt++ {
		resp, err := c.HTTP.Do(req)
		if err != nil {
			return nil, err
		}

		// Retry rate-limited requests, waiting as long as the API asks.
		if resp.StatusCode == http.StatusTooManyRequests && attempt < maxRateLimitRetries {
			wait := retryAfter(resp.Header)
			resp.Body.Close()

			// The body was consumed by the attempt just made, so it has to be
			// rebuilt — otherwise the retry would send an empty body.
			retryReq, err := rewind(req)
			if err != nil {
				return nil, err
			}
			if err := sleep(req.Context(), wait); err != nil {
				return nil, err
			}
			req = retryReq
			continue
		}

		if resp.StatusCode >= 400 {
			body, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr != nil {
				return nil, fmt.Errorf("reading error response (HTTP %d): %w", resp.StatusCode, readErr)
			}
			return nil, parseError(resp.StatusCode, body)
		}

		return resp, nil
	}
}

// rewind returns a request whose body can be sent again. Requests built from an
// in-memory body carry GetBody, which is what makes a retry possible.
func rewind(req *http.Request) (*http.Request, error) {
	if req.Body == nil || req.Body == http.NoBody {
		return req, nil
	}
	if req.GetBody == nil {
		return nil, fmt.Errorf("cannot retry rate-limited request: body is not replayable")
	}

	body, err := req.GetBody()
	if err != nil {
		return nil, fmt.Errorf("rewinding request body: %w", err)
	}
	retryReq := req.Clone(req.Context())
	retryReq.Body = body
	return retryReq, nil
}

// retryAfter reads the wait hint from a 429, preferring Bokio's header over the
// standard one, and clamps it so a bad value cannot stall the process.
func retryAfter(h http.Header) time.Duration {
	for _, name := range []string{"Bokio-RateLimit-RetryAfter", "Retry-After"} {
		value := strings.TrimSpace(h.Get(name))
		if value == "" {
			continue
		}
		seconds, err := strconv.Atoi(value)
		if err != nil || seconds <= 0 {
			continue
		}
		if seconds > maxRetryWaitSeconds {
			seconds = maxRetryWaitSeconds
		}
		return time.Duration(seconds) * time.Second
	}
	return defaultRetryWait
}

// sleep waits for d, or returns early if the context is cancelled.
func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(ctx context.Context, url string, body any) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(ctx context.Context, url string, body any) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, url string) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// PostRaw performs a POST with a raw request (for multipart).
func (c *Client) PostRaw(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	return c.do(req)
}

// GetJSON performs a GET and decodes the JSON response.
func (c *Client) GetJSON(ctx context.Context, url string, v any) error {
	resp, err := c.Get(ctx, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

// PostJSON performs a POST and decodes the JSON response.
func (c *Client) PostJSON(ctx context.Context, url string, body any, v any) error {
	resp, err := c.Post(ctx, url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

// PutJSON performs a PUT and decodes the JSON response.
func (c *Client) PutJSON(ctx context.Context, url string, body any, v any) error {
	resp, err := c.Put(ctx, url, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

// PostEmpty performs a POST expecting no response body.
func (c *Client) PostEmpty(ctx context.Context, url string, body any) error {
	resp, err := c.Post(ctx, url, body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// DeleteEmpty performs a DELETE expecting no response body.
func (c *Client) DeleteEmpty(ctx context.Context, url string) error {
	resp, err := c.Delete(ctx, url)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// BuildListURL appends pagination and query parameters to a base URL.
func BuildListURL(base string, params ListParams) string {
	u, _ := url.Parse(base)
	q := u.Query()
	if params.Page > 0 {
		q.Set("page", fmt.Sprintf("%d", params.Page))
	}
	if params.PageSize > 0 {
		q.Set("pageSize", fmt.Sprintf("%d", params.PageSize))
	}
	if params.Query != "" {
		q.Set("query", params.Query)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
