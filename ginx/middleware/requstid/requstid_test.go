package requstid

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// GetLocalHTTPAddrPairs get available http server and request address
func GetLocalHTTPAddrPairs() (serverAddr string, requestAddr string) {
	port, err := GetAvailablePort()
	if err != nil {
		fmt.Printf("GetAvailablePort error: %v\n", err)
		return "", ""
	}
	serverAddr = fmt.Sprintf(":%d", port)
	requestAddr = fmt.Sprintf("http://127.0.0.1:%d", port)
	return serverAddr, requestAddr
}

// GetAvailablePort get available port
func GetAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()

	return port, err
}

func runRequestIDHTTPServer(fn func(c *gin.Context)) string {
	serverAddr, requestAddr := GetLocalHTTPAddrPairs()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(RequestID())
	r.GET("/ping", func(c *gin.Context) {
		fn(c)
		c.String(200, "pong")
	})

	go func() {
		err := r.Run(serverAddr)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)
	return requestAddr
}

func TestFieldRequestIDFromContext(t *testing.T) {
	requestAddr := runRequestIDHTTPServer(func(c *gin.Context) {
		str := GetCtxRequestID(c)
		t.Log(str)

		str = HeaderRequestID(c)
		t.Log(str)

		str = CtxRequestID(c)
		t.Log(str)

	})

	_, err := http.Get(requestAddr + "/ping")
	assert.NoError(t, err)
}

func TestGetRequestIDFromContext(t *testing.T) {
	str := GetCtxRequestID(&gin.Context{})
	assert.Equal(t, "", str)
	str = CtxRequestID(context.Background())
	assert.Equal(t, "", str)
}

func TestRequestId(t *testing.T) {
	testCases := []struct {
		name           string
		requestBuilder func() *http.Request

		validFunc func(value string) bool
	}{
		{
			name: "Header 里没有 X-Request-ID 参数",
			requestBuilder: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				return req
			},
			validFunc: func(value string) bool {
				return value == ""
			},
		},
		{
			name: "Header 里有 X-Request-ID 参数",
			requestBuilder: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set(HeaderXRequestIDKey, "moocss")
				return req
			},
			validFunc: func(value string) bool {
				return value == "moocss"
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建 gin engine
			engine := gin.Default()
			engine.Use(RequestID())
			engine.GET("/test", func(ctx *gin.Context) {
				ctx.String(200, "test")
			})

			// 创建 request 请求
			req := tc.requestBuilder()
			w := httptest.NewRecorder()
			// 接口调用
			engine.ServeHTTP(w, req)
			assert.True(t, tc.validFunc(w.Header().Get(HeaderXRequestIDKey)))
		})
	}
}
