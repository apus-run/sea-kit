package zlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	logger := NewLogger(WithEncoding("json"), WithFilename("test.log"))
	defer logger.Close()

	logger.Info("This is an info message")
	logger.Error("This is an error message")
}
