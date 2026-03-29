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
	"time"
)

const (
	DefaultBaseURL = "https://api.bokio.se/v1"
)

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
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}

	// Handle rate limiting
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Bokio-RateLimit-RetryAfter")
		resp.Body.Close()
		seconds, _ := strconv.Atoi(retryAfter)
		if seconds <= 0 {
			seconds = 5
		}
		time.Sleep(time.Duration(seconds) * time.Second)
		return c.HTTP.Do(req)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var errResp errorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Code != "" {
			return nil, &errResp.Error
		}
		return nil, fmt.Errorf("API error: %d %s", resp.StatusCode, string(body))
	}

	return resp, nil
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
