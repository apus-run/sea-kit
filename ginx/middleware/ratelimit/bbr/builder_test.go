package bbr

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	ratelimit "github.com/apus-run/sea-kit/ratelimit_bbr"
)

type ratelimitMock struct {
	reached bool
}

func (r *ratelimitMock) Allow() (ratelimit.DoneFunc, error) {
	return func(_ ratelimit.DoneInfo) {
		r.reached = true
	}, nil
}

func TestServer(t *testing.T) {
	rlm := &ratelimitMock{}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(NewBuilder().Limiter(rlm).Build())

	if !rlm.reached {
		t.Log("The ratelimit must run the done function.")
	}

	// Add a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test route")
	})

	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedBody := "Test route"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected response body %q, but got %q", expectedBody, rr.Body.String())
	}
}
