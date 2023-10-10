package accesslog

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"

	"github.com/apus-run/sea-kit/ginx/middleware/requstid"
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
	// Create a slog logger, which:
	//   - Logs to stdout.
	w := os.Stdout
	logger := slog.New(
		slog.NewJSONHandler(
			w,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		),
	)
	logger.WithGroup("http").
		With("environment", "production").
		With("server", "gin/1.9.0").
		With("server_start_time", time.Now()).
		With("gin_mode", gin.EnvGinMode)
	// [SetDefault]还更新了[log]包使用的默认logger
	slog.SetDefault(logger)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	engine.Use(requstid.RequestID())
	// custom print log
	engine.Use(NewBuilder(
		func(ctx context.Context, al *AccessLog) {
			logger.Debug("Gin 收到请求", slog.Any("req", al))
		}).
		AllowReqBody().
		AllowRespBody().
		MaxLength(1024).
		IgnoreRoutes("/foo").
		Build())

	handler := GinHandler(engine)

	// run server using httptest
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
