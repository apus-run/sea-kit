package requstid

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/apus-run/sea-kit/utils"
)

func runRequestIDHTTPServer(fn func(c *gin.Context)) string {
	serverAddr, requestAddr := utils.GetLocalHTTPAddrPairs()

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
