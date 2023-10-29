package zlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	logger := NewLogger(WithEncoding("json"), WithFilename("test.log"))
	defer logger.Close()

	logger.Info("This is an info message")
	logger.Infof("我是日志: %v, %v", String("route", "/hello"), Int64("port", 8090))
	logger.Error("This is an error message")
}
