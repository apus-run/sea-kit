package zlog

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestGlobalLogInit(t *testing.T) {
	la := &loggerAppliance{}
	logger := NewLogger(WithEncoding("json"), WithFilename("test.log"))
	la.SetLogger(logger)

	la.Info("test info")
	la.Error("test error")
	la.Warn("test warn")

	la.Error("处理业务逻辑出错",
		zap.String("path", "/global.go"),
		// 命中的路由
		zap.String("route", "/hello"),
		zap.Error(errors.New("自定义错误")),
		zap.Time("time", time.Now()),
		zap.Duration("duration", time.Duration(int64(10))),
	)

	la.Info("test info")
}

func TestGlobalLog(t *testing.T) {
	L().Info("test info")
	L().Error("test error")
	L().Warn("test warn")

	L().Error("处理业务逻辑出错",
		zap.String("path", "/global.go"),
		// 命中的路由
		zap.String("route", "/hello"),
		zap.Error(errors.New("自定义错误")),
		zap.Time("time", time.Now()),
		zap.Duration("duration", time.Duration(int64(10))),
	)
}

func TestContext(t *testing.T) {
	c := WithContext(context.Background(), "name", "foo")

	C(c).Info("test info")
	C(c).Error("test error")
	C(c).Warn("test warn")
	C(c).Error("处理业务逻辑出错",
		zap.String("path", "/global.go"),
		// 命中的路由
		zap.String("route", "/hello"),
		zap.Error(errors.New("自定义错误")),
		zap.Time("time", time.Now()),
		zap.Duration("duration", time.Duration(int64(10))),
	)
}
