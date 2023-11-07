package slogx

import (
	"context"
	"log"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"
)

type Password string

func (Password) LogValue() slog.Value {
	return slog.StringValue("******")
}

func TestLog(t *testing.T) {
	logger := NewLogger(WithEncoding("json"), WithFilename("test.log"))
	logger.Debug("This is a debug message", slog.Any("key", "value"))
	logger.Info("This is a info message")
	logger.Warn("This is a warn message")
	logger.Error("This is a error message")

	logger.Info("WebServer服务信息",
		slog.Group("http",
			slog.Int("status", 200),
			slog.String("method", "POST"),
			slog.Time("time", time.Now()),
		),
	)

	log.Print("This is a print message")

	slog.Info("敏感数据", slog.Any("password", Password("1234567890")))

	// 设置slog为默认日志
	slog.SetDefault(logger)
}

func TestNewLogger(t *testing.T) {
	var called int32
	ctx := context.WithValue(context.Background(), "foobar", "helloworld")

	logger := NewLogger(
		WithEncoding("json"),
		WithLogLevel("debug"),
		WithFilename("test.log"),
	)
	ApplyHandlerOption(WithHandleFunc(func(ctx context.Context, r *slog.Record) {
		r.AddAttrs(slog.String("value", ctx.Value("foobar").(string)))
		atomic.AddInt32(&called, 1)
	}))

	//logger = logger.With(slog.String("sub_logger", "true"))
	//ctx = NewContext(ctx, logger)
	//logger = FromContext(ctx)
	//logger.InfoContext(ctx, "print something")

	logger = logger.With(slog.String("sub_logger", "true"))
	ctx = WithContext(ctx, logger)
	logger = FromContext(ctx)
	t.Logf("%#v", logger)
	logger.InfoContext(ctx, "print something")

	if atomic.LoadInt32(&called) != 1 {
		t.FailNow()
	}
}
