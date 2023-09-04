package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewCORS(t *testing.T) {
	cors := NewCORS(
		WithAllowOrigins([]string{"http://www.baidu.com"}),
		WithAllowMethods([]string{http.MethodGet}),
		WithAllowCredentials(true),
		WithAllowHeaders([]string{"Content-Type", "Authorization"}),
		WithExposeHeaders([]string{"Range"}),
		WithMaxAge(24*time.Hour),
	)
	assert.Equal(t, []string{"http://www.baidu.com"}, cors.allowOrigins)
	assert.Equal(t, []string{http.MethodGet}, cors.allowMethods)
	assert.Equal(t, true, cors.allowCredentials)
	assert.Equal(t, []string{"Content-Type", "Authorization"}, cors.allowHeaders)
	assert.Equal(t, []string{"Range"}, cors.exposeHeaders)
	assert.Equal(t, 24*time.Hour, cors.maxAge)

}

func TestCORS(t *testing.T) {
	tests := []struct {
		name      string
		origins   []string
		reqOrigin string
		expect    string
	}{
		{
			name:   "allow all origins",
			expect: "*",
		},
		{
			name:      "allow one origin",
			origins:   []string{"http://local"},
			reqOrigin: "http://local",
			expect:    "http://local",
		},
		{
			name:      "allow many origins",
			origins:   []string{"http://local", "http://remote"},
			reqOrigin: "http://local",
			expect:    "http://local",
		},
		{
			name:      "allow all origins",
			reqOrigin: "http://local",
			expect:    "*",
		},
		{
			name:      "allow many origins with all mark",
			origins:   []string{"http://local", "http://remote", "*"},
			reqOrigin: "http://another",
			expect:    "*",
		},
		{
			name:      "not allow origin",
			origins:   []string{"http://local", "http://remote"},
			reqOrigin: "http://another",
		},
	}

	methods := []string{
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
	}

	for _, test := range tests {
		for _, method := range methods {
			test := test
			t.Run(test.name+"-handler", func(t *testing.T) {
				r := httptest.NewRequest(method, "http://localhost", nil)
				r.Header.Set(originHeader, test.reqOrigin)
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				ctx.Request = r
				if len(test.origins) != 0 {
					NewCORS(WithAllowOrigins(test.origins)).Build()(ctx)
				} else {
					NewCORS().Build()(ctx)
				}

				if method == http.MethodOptions {
					assert.Equal(t, http.StatusNoContent, ctx.Writer.Status())
				} else {
					assert.Equal(t, http.StatusOK, ctx.Writer.Status())
				}
				assert.Equal(t, test.expect, ctx.Writer.Header().Get(allowOrigin))
			})
			t.Run(test.name+"-handler-custom", func(t *testing.T) {
				r := httptest.NewRequest(method, "http://localhost", nil)
				r.Header.Set(originHeader, test.reqOrigin)
				w := httptest.NewRecorder()
				w.Header().Set("foo", "bar")
				ctx, _ := gin.CreateTestContext(w)
				ctx.Request = r
				if len(test.origins) != 0 {
					NewCORS(WithAllowOrigins(test.origins)).Build()(ctx)
				} else {
					NewCORS().Build()(ctx)
				}
				if method == http.MethodOptions {
					assert.Equal(t, http.StatusNoContent, ctx.Writer.Status())
				} else {
					assert.Equal(t, http.StatusOK, ctx.Writer.Status())
				}
				assert.Equal(t, test.expect, ctx.Writer.Header().Get(allowOrigin))
				assert.Equal(t, "bar", ctx.Writer.Header().Get("foo"))
			})
		}
	}

	for _, test := range tests {
		for _, method := range methods {
			test := test
			t.Run(test.name+"-middleware", func(t *testing.T) {
				r := httptest.NewRequest(method, "http://localhost", nil)
				r.Header.Set(originHeader, test.reqOrigin)
				w := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(w)
				ctx.Request = r
				if len(test.origins) != 0 {
					NewCORS(WithAllowOrigins(test.origins)).Build()(ctx)
				} else {
					NewCORS().Build()(ctx)
				}
				if method == http.MethodOptions {
					assert.Equal(t, http.StatusNoContent, ctx.Writer.Status())
				} else {
					assert.Equal(t, http.StatusOK, ctx.Writer.Status())
				}
				assert.Equal(t, test.expect, ctx.Writer.Header().Get(allowOrigin))
			})
			t.Run(test.name+"-middleware-custom", func(t *testing.T) {
				r := httptest.NewRequest(method, "http://localhost", nil)
				r.Header.Set(originHeader, test.reqOrigin)
				w := httptest.NewRecorder()
				w.Header().Set("foo", "bar")
				ctx, _ := gin.CreateTestContext(w)
				ctx.Request = r
				if len(test.origins) != 0 {
					NewCORS(WithAllowOrigins(test.origins)).Build()(ctx)
				} else {
					NewCORS().Build()(ctx)
				}
				if method == http.MethodOptions {
					assert.Equal(t, http.StatusNoContent, ctx.Writer.Status())
				} else {
					assert.Equal(t, http.StatusOK, ctx.Writer.Status())
				}
				assert.Equal(t, test.expect, ctx.Writer.Header().Get(allowOrigin))
				assert.Equal(t, "bar", ctx.Writer.Header().Get("foo"))
			})
		}
	}
}
