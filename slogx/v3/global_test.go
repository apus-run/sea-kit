package slogx

import (
	"testing"
)

func TestGlobalLog(t *testing.T) {

	logger := NewLogger()
	SetLogger(logger)

	Info("test info")
	Error("test error")
	Warn("test warn")
}
