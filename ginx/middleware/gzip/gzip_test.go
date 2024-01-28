package gzip

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
)

func GinHandler(r *gin.Engine) *gin.Engine {
	helloFun := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "hello world",
		})
	}

	pingFun := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ping",
		})
	}

	fooFun := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "foo",
		})
	}

	r.GET("/foo", fooFun)
	r.GET("/hello", helloFun)
	r.GET("/ping", pingFun)
	r.DELETE("/hello", helloFun)
	r.POST("/hello", helloFun)
	r.PUT("/hello", helloFun)
	r.PATCH("/hello", helloFun)

	return r
}

func TestRequest(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	engine.Use()

	handler := GinHandler(engine)

	engine.NoRoute(NewBuilder().
		SetLevel(DefaultCompression).
		Build(WithGzipExcludedExtensions([]string{"", ".html"})), func(c *gin.Context) {

		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d, public", 31536000))
		c.Header("Expires", time.Now().AddDate(1, 0, 0).Format("Mon, 01 Jan 2006 00:00:00 GMT"))

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "404",
		})
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/ping").
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("msg", "ping")
	e.GET("/foo").
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("msg", "foo")
	e.GET("/hello").
		Expect().
		Status(http.StatusOK).JSON().Object().HasValue("msg", "hello world")
}
