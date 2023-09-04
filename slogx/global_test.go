package slogx

import (
	"testing"
)

func TestGlobalLog(t *testing.T) {
	la := &loggerAppliance{}
	logger := NewLogger()
	la.SetLogger(*logger)

	Info("test info")
	Error("test error")
	Warn("test warn")
}
