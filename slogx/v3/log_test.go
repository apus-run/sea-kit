package slogx

import (
	"context"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Password string

func (Password) LogValue() slog.Value {
	return slog.StringValue("******")
}

func TestLog(t *testing.T) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "logs.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   true,
	}

	logger := NewLogger(WithFormat(FormatJSON), WithWriter(lumberJackLogger))
	logger.Debug("This is a debug message", Any("key", "value"))
	logger.Info("This is a info message")
	logger.Warn("This is a warn message")
	logger.Error("This is a error message")
	logger.Fatalf("ssssss%v", "22222")

	logger.Info("WebServer服务信息",
		Group("http",
			Int("status", 200),
			String("method", "POST"),
			Time("time", time.Now()),
		),
	)

	logger.Info("敏感数据", Any("password", Password("1234567890")))
}

func TestNewLogger(t *testing.T) {
	var called int32
	ctx := context.WithValue(context.Background(), "foobar", "helloworld")

	logger := NewLogger(
		WithFormat(FormatJSON),
		WithLogLevel("debug"),
	)
	ApplyHandlerOption(WithHandleFunc(func(ctx context.Context, r *slog.Record) {
		r.AddAttrs(String("value", ctx.Value("foobar").(string)))
		atomic.AddInt32(&called, 1)
	}))

	// logger = With(String("sub_logger", "true"))
	// ctx = NewContext(ctx, logger)
	// logger = FromContext(ctx)
	// logger.InfoContext(ctx, "print something")
	l := logger.WithGroup("moocss").With(String("sub_logger", "true"))
	ctx = WithContext(ctx, l)
	logger = FromContext(ctx)
	t.Logf("%#v", logger)
	logger.InfoContext(ctx, "print something")

	if atomic.LoadInt32(&called) != 1 {
		t.FailNow()
	}
}
