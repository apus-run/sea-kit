package slogx

import (
	"log"
	"log/slog"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	l := NewLogger(WithEncoding("json"), WithFilename("test.log"))
	l.Debug("This is a debug message", slog.Any("key", "value"))
	l.Info("This is a info message")
	l.Warn("This is a warn message")
	l.Error("This is a error message")

	l.Info("WebServer服务信息",
		slog.Group("http",
			slog.Int("status", 200),
			slog.String("method", "POST"),
			slog.Time("time", time.Now()),
		),
	)

	log.Print("This is a print message")

	slog.SetDefault(l.logger)
}
