package zlog

import (
	"sync"

	"go.uber.org/zap"
)

// globalLogger is designed as a global logger in current process.
var global = &loggerAppliance{}

// loggerAppliance is the proxy of `Logger` to
// make logger change will affect all sub-logger.
type loggerAppliance struct {
	lock sync.Mutex
	Logger
}

func init() {
	logger := NewLogger()

	global.SetLogger(*logger)
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

// SetLogger should be called before any other log call.
// And it is NOT THREAD SAFE.
func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

// GetLogger returns global logger appliance as logger in current process.
func GetLogger() Logger {
	return global.Logger
}

// Info uses the fmt.Sprint to construct and log a message using the DefaultLogger.
func Info(msg string, tags ...zap.Field) {
	global.Info(msg, tags...)
}

// Infof logs a message at info level.
func Infof(format string, args ...any) {
	global.Infof(format, args...)
}

// Error uses the fmt.Sprint to construct and log a message using the DefaultLogger.
func Error(msg string, tags ...zap.Field) {
	global.Error(msg, tags...)
}

// Errorf logs a message at error level.
func Errorf(format string, args ...any) {
	global.Errorf(format, args...)
}

// Debug uses the fmt.Sprint to construct and log a message using the DefaultLogger.
func Debug(msg string, tags ...zap.Field) {
	global.Debug(msg, tags...)
}

// Debugf logs a message at debug level.
func Debugf(format string, args ...any) {
	global.Debugf(format, args...)
}

// Warn uses the fmt.Sprint to construct and log a message using the DefaultLogger.
func Warn(msg string, tags ...zap.Field) {
	global.Warn(msg, tags...)
}

// Warnf logs a message at warn level.
func Warnf(format string, args ...any) {
	global.Warnf(format, args...)
}

// Fatal uses the fmt.Sprint to construct and log a message using the DefaultLogger.
func Fatal(msg string, tags ...zap.Field) {
	global.Fatal(msg, tags...)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, args ...any) {
	global.Fatalf(format, args...)
}

// Panic uses the fmt.Sprint to construct and log a message using the DefaultLogger.
func Panic(msg string, tags ...zap.Field) {
	global.Panic(msg, tags...)
}

// Panicf logs a message at panic level.
func Panicf(format string, args ...any) {
	global.Panicf(format, args...)
}
