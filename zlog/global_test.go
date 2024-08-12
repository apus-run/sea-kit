package zlog

import (
	"errors"
	"testing"
	"time"
)

func TestGlobalLog(t *testing.T) {
	la := &loggerAppliance{}
	logger := NewLogger()
	la.SetLogger(logger)

	la.Info("test info")
	la.Error("test error")
	la.Warn("test warn")

	la.Error("处理业务逻辑出错",
		String("path", "/global.go"),
		// 命中的路由
		String("route", "/hello"),
		Error(errors.New("自定义错误")),
		Time("time", time.Now()),
		Duration("duration", time.Duration(int64(10))),
	)

}
