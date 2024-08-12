// 源码来自: github.com/davesavic/clink

package got

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	HttpClient      *http.Client
	Headers         map[string]string
	RateLimiter     *rate.Limiter
	ReadLimitBytes  int64
	RetryMax        int
	ShouldRetryFunc func(*http.Request, *http.Response, error) bool
}

// NewClient creates a new client with the given options.
func NewClient(opts ...Option) *Client {
	return Apply(opts...)
}

// Do sends the given request and returns the response.
// If the request is rate limited, the client will wait for the rate limiter to allow the request.
// If the request fails, the client will retry the request the number of times specified by MaxRetries.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	if c.RateLimiter != nil {
		if err := c.RateLimiter.Wait(req.Context()); err != nil {
			return nil, fmt.Errorf("failed to wait for rate limiter: %w", err)
		}
	}

	var resp *http.Response
	var body []byte
	var err error

	if req.Body != nil && req.Body != http.NoBody {
		body, err = io.ReadAll(io.LimitReader(req.Body, c.ReadLimitBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}

		err = req.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close request body: %w", err)
		}
	}

	for attempt := 0; attempt <= c.RetryMax; attempt++ {
		if len(body) > 0 {
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		resp, err = c.HttpClient.Do(req)

		if req.Context().Err() != nil {
			return nil, fmt.Errorf("request context error: %w", req.Context().Err())
		}

		if c.ShouldRetryFunc != nil && !c.ShouldRetryFunc(req, resp, err) {
			break
		}

		if attempt < c.RetryMax {
			select {
			case <-time.After(time.Duration(attempt) * time.Second):
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	return resp, nil
}

// Head sends a HEAD request to the given URL.
func (c *Client) Head(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Get sends a GET request to the given URL.
func (c *Client) Options(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodOptions, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Get sends a GET request to the given URL.
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post sends a POST request to the given URL with the given body.
func (c *Client) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Put sends a PUT request to the given URL.
func (c *Client) Put(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Patch sends a PATCH request to the given URL.
func (c *Client) Patch(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Delete sends a DELETE request to the given URL.
func (c *Client) Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// ResponseToJson decodes the response body into the target.
func ResponseToJson(response *http.Response, target any) error {
	if response == nil {
		return fmt.Errorf("response is nil")
	}

	if response.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func GetCompleteURL(originURL string, params map[string]string) string {
	values := neturl.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	queriesStr, _ := neturl.QueryUnescape(values.Encode())
	if len(queriesStr) == 0 {
		return originURL
	}
	return fmt.Sprintf("%s?%s", originURL, queriesStr)
}
