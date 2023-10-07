package metrics

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestMetrics(t *testing.T) {
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

	engine.Use(Metrics(engine,
		WithMetricsPath("/metrics"),
		WithIgnoreStatusCodes(http.StatusNotFound),
		WithIgnoreRequestPaths("/hello-ignore"),
		WithIgnoreRequestMethods(http.MethodDelete),
	))

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

	str := e.GET("/metrics").Expect().Status(http.StatusOK).Text()
	logger.Info("输出值: %v", slog.Any("text:", str))
}
