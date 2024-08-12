package got

import (
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Option is config option.
type Option func(*Client)

func defaultClient() *Client {
	return &Client{
		HttpClient:     http.DefaultClient,
		Headers:        make(map[string]string),
		ReadLimitBytes: 4 * 1024 * 1024, // 单次读取限制 4M
	}
}

func Apply(opts ...Option) *Client {
	c := defaultClient()

	for _, o := range opts {
		o(c)
	}
	return c
}

// WithClient .
func WithClient(client *http.Client) Option {
	return func(o *Client) {
		o.HttpClient = client
	}
}

// WithHeader sets a header for the client.
func WithHeader(key, value string) Option {
	return func(c *Client) {
		c.Headers[key] = value
	}
}

// WithHeaders sets the headers for the client.
func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		for key, value := range headers {
			c.Headers[key] = value
		}
	}
}

func WithReadLimitBytes(limit int64) Option {
	return func(j *Client) {
		j.ReadLimitBytes = limit
	}
}

// WithRateLimit sets the rate limit for the client in requests per minute.
func WithRateLimit(rpm int) Option {
	return func(c *Client) {
		interval := time.Minute / time.Duration(rpm)
		c.RateLimiter = rate.NewLimiter(rate.Every(interval), 1)
	}
}

// WithBasicAuth sets the basic auth header for the client.
func WithBasicAuth(username, password string) Option {
	return func(c *Client) {
		auth := username + ":" + password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		c.Headers["Authorization"] = "Basic " + encodedAuth
	}
}

// WithBearerAuth sets the bearer auth header for the client.
func WithBearerAuth(token string) Option {
	return func(c *Client) {
		c.Headers["Authorization"] = "Bearer " + token
	}
}

// WithUserAgent sets the user agent header for the client.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		c.Headers["User-Agent"] = ua
	}
}

// WithRetries sets the retry count and retry function for the client.
func WithRetries(count int, retryFunc func(*http.Request, *http.Response, error) bool) Option {
	return func(c *Client) {
		c.RetryMax = count
		c.ShouldRetryFunc = retryFunc
	}
}
