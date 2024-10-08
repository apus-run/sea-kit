package got

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	testCases := []struct {
		name   string
		opts   []Option
		result func(*Client) bool
	}{
		{
			name: "default client with no options",
			opts: []Option{},
			result: func(client *Client) bool {
				return client.HttpClient != nil && client.Headers != nil && len(client.Headers) == 0
			},
		},
		{
			name: "client with custom http client",
			opts: []Option{
				WithClient(nil),
			},
			result: func(client *Client) bool {
				return client.HttpClient == nil
			},
		},
		{
			name: "client with custom headers",
			opts: []Option{
				WithHeaders(map[string]string{"key": "value"}),
			},
			result: func(client *Client) bool {
				return client.Headers != nil && len(client.Headers) == 1
			},
		},
		{
			name: "client with custom header",
			opts: []Option{
				WithHeader("key", "value"),
			},
			result: func(client *Client) bool {
				return client.Headers != nil && len(client.Headers) == 1
			},
		},
		{
			name: "client with custom rate limit",
			opts: []Option{
				WithRateLimit(60),
			},
			result: func(client *Client) bool {
				return client.RateLimiter != nil && client.RateLimiter.Limit() == 1
			},
		},
		{
			name: "client with basic auth",
			opts: []Option{
				WithBasicAuth("username", "password"),
			},
			result: func(client *Client) bool {
				b64, err := base64.StdEncoding.DecodeString(
					strings.Replace(client.Headers["Authorization"], "Basic ", "", 1),
				)
				if err != nil {
					return false
				}

				return string(b64) == "username:password"
			},
		},
		{
			name: "client with bearer token",
			opts: []Option{
				WithBearerAuth("token"),
			},
			result: func(client *Client) bool {
				return client.Headers["Authorization"] == "Bearer token"
			},
		},
		{
			name: "client with user agent",
			opts: []Option{
				WithUserAgent("user-agent"),
			},
			result: func(client *Client) bool {
				return client.Headers["User-Agent"] == "user-agent"
			},
		},
		{
			name: "client with retries",
			opts: []Option{
				WithRetries(3, func(request *http.Request, response *http.Response, err error) bool {
					return true
				}),
			},
			result: func(client *Client) bool {
				return client.RetryMax == 3 && client.ShouldRetryFunc != nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewClient(tc.opts...)

			if c == nil {
				t.Error("expected client to be created")
			}

			if !tc.result(c) {
				t.Errorf("expected client to be created with options: %+v", tc.opts)
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	testCases := []struct {
		name        string
		opts        []Option
		setupServer func() *httptest.Server
		resultFunc  func(*http.Response, error) bool
	}{
		{
			name: "successful response no body",
			opts: []Option{},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			resultFunc: func(response *http.Response, err error) bool {
				return response != nil && err == nil && response.StatusCode == http.StatusOK
			},
		},
		{
			name: "successful response with text body",
			opts: []Option{},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte("response"))
				}))
			},
			resultFunc: func(response *http.Response, err error) bool {
				bodyContents, err := io.ReadAll(response.Body)
				if err != nil {
					return false
				}

				return string(bodyContents) == "response"
			},
		},
		{
			name: "successful response with json body",
			opts: []Option{},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_ = json.NewEncoder(w).Encode(map[string]string{"key": "value"})
				}))
			},
			resultFunc: func(response *http.Response, err error) bool {
				var target map[string]string
				er := ResponseToJson(response, &target)
				if er != nil {
					return false
				}

				return target["key"] == "value"
			},
		},
		{
			name: "successful response with json body and custom headers",
			opts: []Option{
				WithHeaders(map[string]string{"key": "value"}),
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Header.Get("key") != "value" {
						w.WriteHeader(http.StatusBadRequest)
					}

					_ = json.NewEncoder(w).Encode(map[string]string{"key": "value"})
				}))
			},
			resultFunc: func(response *http.Response, err error) bool {
				var target map[string]string
				er := ResponseToJson(response, &target)
				if er != nil {
					return false
				}

				return target["key"] == "value"
			},
		},
		{
			name: "successful response with json body and custom header",
			opts: []Option{
				WithHeader("key", "value"),
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Header.Get("key") != "value" {
						w.WriteHeader(http.StatusBadRequest)
					}

					_ = json.NewEncoder(w).Encode(map[string]string{"key": "value"})
				}))
			},
			resultFunc: func(response *http.Response, err error) bool {
				var target map[string]string
				er := ResponseToJson(response, &target)
				if er != nil {
					return false
				}

				return target["key"] == "value"
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := tc.setupServer()
			defer server.Close()

			opts := append(tc.opts, WithClient(server.Client()))
			c := NewClient(opts...)

			if c == nil {
				t.Error("expected client to be created")
			}

			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			if err != nil {
				t.Errorf("failed to create request: %v", err)
			}

			resp, err := c.Do(req)
			if !tc.resultFunc(resp, err) {
				t.Errorf("expected result to be successful")
			}
		})
	}
}

func TestClient_Methods(t *testing.T) {
	serverFunc := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Method", r.Method)
		}))
	}
	resultFunc := func(r *http.Response, m string) bool {
		return r.Header.Get("X-Method") == m
	}
	testCases := []struct {
		name        string
		method      string
		body        io.Reader
		setupServer func() *httptest.Server
		resultFunc  func(*http.Response, string) bool
	}{
		{
			name:        "successful head response",
			method:      http.MethodHead,
			setupServer: serverFunc,
			resultFunc:  resultFunc,
		},
		{
			name:        "successful options response",
			method:      http.MethodOptions,
			setupServer: serverFunc,
			resultFunc:  resultFunc,
		},
		{
			name:        "successful get response",
			method:      http.MethodGet,
			setupServer: serverFunc,
			resultFunc:  resultFunc,
		},
		{
			name:        "successful post response",
			method:      http.MethodPost,
			setupServer: serverFunc,
			resultFunc:  resultFunc,
		},
		{
			name:        "successful put response",
			method:      http.MethodPut,
			setupServer: serverFunc,
			resultFunc:  resultFunc,
		},
		{
			name:        "successful patch response",
			method:      http.MethodPatch,
			setupServer: serverFunc,
			resultFunc:  resultFunc,
		},
		{
			name:        "successful delete response",
			method:      http.MethodDelete,
			setupServer: serverFunc,
			resultFunc:  resultFunc,
		},
	}

	call := func(c *Client, method, url string, body io.Reader) (*http.Response, error) {
		switch method {
		case http.MethodHead:
			return c.Head(url)
		case http.MethodOptions:
			return c.Options(url)
		case http.MethodGet:
			return c.Get(url)
		case http.MethodPost:
			return c.Post(url, body)
		case http.MethodPut:
			return c.Put(url, body)
		case http.MethodPatch:
			return c.Patch(url, body)
		case http.MethodDelete:
			return c.Delete(url)
		}
		return nil, nil
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := tc.setupServer()
			defer server.Close()
			c := NewClient(WithClient(server.Client()))
			if c == nil {
				t.Error("expected client to be created")
			}
			resp, _ := call(c, tc.method, server.URL, tc.body)
			if !tc.resultFunc(resp, tc.method) {
				t.Errorf("expected result to be successful")
			}
		})
	}
}

func TestClient_ResponseToJson(t *testing.T) {
	testCases := []struct {
		name       string
		response   *http.Response
		target     any
		resultFunc func(*http.Response, any) bool
	}{
		{
			name: "successful response with json body",
			response: &http.Response{
				Body: io.NopCloser(strings.NewReader(`{"key": "value"}`)),
			},
			resultFunc: func(response *http.Response, target any) bool {
				var t map[string]string
				er := ResponseToJson(response, &t)
				if er != nil {
					return false
				}

				return t["key"] == "value"
			},
		},
		{
			name:     "response is nil",
			response: nil,
			resultFunc: func(response *http.Response, target any) bool {
				var t map[string]string
				er := ResponseToJson(response, &t)
				if er == nil {
					return false
				}

				return er.Error() == "response is nil"
			},
		},
		{
			name: "response body is nil",
			response: &http.Response{
				Body: nil,
			},
			resultFunc: func(response *http.Response, target any) bool {
				var t map[string]string
				er := ResponseToJson(response, &t)
				if er == nil {
					return false
				}

				return er.Error() == "response body is nil"
			},
		},
		{
			name: "json decode error",
			response: &http.Response{
				Body: io.NopCloser(strings.NewReader(`{"key": "value`)),
			},
			target: nil,
			resultFunc: func(response *http.Response, target any) bool {
				var t map[string]string
				er := ResponseToJson(response, &t)
				if er == nil {
					return false
				}

				return strings.Contains(er.Error(), "failed to decode response")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.resultFunc(tc.response, tc.target) {
				t.Errorf("expected result to be successful")
			}
		})
	}
}

func TestRateLimiter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(
		WithRateLimit(60),
		WithClient(server.Client()),
	)

	startTime := time.Now()

	for i := 0; i < 2; i++ {
		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		if err != nil {
			t.Errorf("failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("failed to make request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code to be 200")
		}
	}

	elapsedTime := time.Since(startTime)
	if elapsedTime.Seconds() < 0.5 || elapsedTime.Seconds() > 1.5 {
		t.Errorf("expected elapsed time to be between 0.5 and 1.5 seconds, got: %f", elapsedTime.Seconds())
	}
}

func TestSuccessfulRetries(t *testing.T) {
	var requestCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++ // Increment the request count
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	retryCount := 3
	client := NewClient(
		WithRetries(retryCount, func(request *http.Request, response *http.Response, err error) bool {
			// Check if the response is a 500 Internal Server Error
			return response != nil && response.StatusCode == http.StatusInternalServerError
		}),
		WithClient(server.Client()),
	)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}

	if requestCount != retryCount+1 { // +1 for the initial request
		t.Errorf("expected %d retries (total requests: %d), but got %d", retryCount, retryCount+1, requestCount)
	}
}

func TestUnsuccessfulRetries(t *testing.T) {
	var requestCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++ // Increment the request count
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	retryCount := 3
	client := NewClient(
		WithRetries(retryCount, func(request *http.Request, response *http.Response, err error) bool {
			return false
		}),
		WithClient(server.Client()),
	)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	_, err = client.Do(req)

	if requestCount != 1 { // +1 for the initial request
		t.Errorf("expected %d retries (total requests: %d), but got %d", retryCount, retryCount+1, requestCount)
	}
}

// TestRequestBodyEmptyOnRetries tests that the request body on a custom io.Reader wrapper is NOT empty on retries.
type oneTimeReaderWrapper struct {
	data     []byte
	consumed bool
}

func (r *oneTimeReaderWrapper) Read(p []byte) (n int, err error) {
	if r.consumed {
		return 0, fmt.Errorf("body already read")
	}
	n = copy(p, r.data)
	r.consumed = true
	return n, io.EOF
}

func TestRequestBodyNotEmptyOnRetries(t *testing.T) {
	var requestCount int
	var lastRequestBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		lastRequestBody = string(bodyBytes)

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(
		WithRetries(1, func(request *http.Request, response *http.Response, err error) bool {
			return true
		}),
		WithClient(server.Client()),
	)

	requestBody := []byte("test body")
	req, err := http.NewRequest(http.MethodPost, server.URL, &oneTimeReaderWrapper{data: requestBody})
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}

	if requestCount != 2 {
		t.Fatalf("expected 2 requests due to retry, but got %d", requestCount)
	}

	expectedBody := string(requestBody)
	if lastRequestBody != expectedBody {
		t.Errorf("expected request body to be '%s' on retry, got '%s'", expectedBody, lastRequestBody)
	}
}

func TestContextCancellationDuringRetries(t *testing.T) {
	var requestCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusInternalServerError) // Always return an error to trigger retries
	}))
	defer server.Close()

	client := NewClient(
		WithRetries(3, func(request *http.Request, response *http.Response, err error) bool {
			// Always return true to retry
			return true
		}),
		WithClient(server.Client()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	_, err = client.Do(req)

	if requestCount > 2 {
		t.Errorf("expected at most 2 requests due to context cancellation, but got %d", requestCount)
	}

	if err == nil || !errors.Is(err, context.Canceled) {
		t.Errorf("expected context cancellation error, but got: %v", err)
	}
}

func TestRequestWithCanceledContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Simulate a delay in the response
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(WithClient(server.Client()))

	ctx, cancel := context.WithCancel(context.Background())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	cancel() // Cancel the context immediately

	_, err = client.Do(req)

	if err == nil || !errors.Is(err, context.Canceled) {
		t.Errorf("expected context cancellation error, but got: %v", err)
	}
}
