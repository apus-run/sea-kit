package slogx

import (
	"errors"
	"log/slog"
	"testing"
	"time"
)

func TestGlobalLog(t *testing.T) {
	la := &loggerAppliance{}
	logger := New()
	la.SetLogger(*logger)

	la.Info("test info")
	la.Error("test error")
	la.Warn("test warn")

	la.Error("处理业务逻辑出错",
		slog.String("path", "/global.go"),
		// 命中的路由
		slog.String("route", "/hello"),
		ErrorString(errors.New("自定义错误")),
		slog.Time("time", time.Now()),
		slog.Duration("duration", time.Duration(int64(10))),
	)
}
