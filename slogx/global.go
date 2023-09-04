package slogx

import (
	"log/slog"
	"sync"
)

// globalLogger is designed as a global logger in current process.
var global = &loggerAppliance{}

// loggerAppliance is the proxy of `Logger` to
// make logger change will affect all sub-logger.
type loggerAppliance struct {
	lock sync.Mutex
	slog.Logger
}

func init() {
	logger := NewLogger()

	global.SetLogger(*logger)
}

func (a *loggerAppliance) SetLogger(in slog.Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

// SetLogger should be called before any other log call.
// And it is NOT THREAD SAFE.
func SetLogger(logger slog.Logger) {
	global.SetLogger(logger)
}

// GetLogger returns global logger appliance as logger in current process.
func GetLogger() slog.Logger {
	return global.Logger
}

func Info(msg string, args ...any) {
	global.Info(msg, args...)
}

func Error(msg string, args ...any) {
	global.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	global.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	global.Warn(msg, args...)
}
