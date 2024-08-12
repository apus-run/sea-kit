package zlog

import (
	"sync"
)

// globalLogger is designed as a global logger in current process.
var global = loggerAppliance{}

// loggerAppliance is the proxy of `Logger` to
// make logger change will affect all sub-logger.
type loggerAppliance struct {
	lock sync.Mutex
	Logger
}

func init() {
	logger := NewLogger()

	global.SetLogger(logger)
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

// L 是 GetLogger 简写
func L() Logger {
	return global.Logger
}

// Default 是 GetLogger
func Default() Logger {
	return global.Logger
}

// Debugf 输出 debug 级别的日志.
func Debugf(format string, args ...any) {
	L().Debugf(format, args...)
}

// Debug 输出 debug 级别的日志.
func Debug(msg string, tags ...Field) {
	L().Debug(msg, tags...)
}

// Infof 输出 info 级别的日志.
func Infof(format string, args ...any) {
	L().Infof(format, args...)
}

// Info 输出 info 级别的日志.
func Info(msg string, tags ...Field) {
	L().Info(msg, tags...)
}

// Warnf 输出 warning 级别的日志.
func Warnf(format string, args ...any) {
	L().Warnf(format, args...)
}

// Warn 输出 warning 级别的日志.
func Warn(msg string, tags ...Field) {
	L().Warn(msg, tags...)
}

// Errorf 输出 error 级别的日志.
func Errorf(format string, args ...any) {
	L().Errorf(format, args...)
}

// Error 输出 error 级别的日志.
func Error(msg string, tags ...Field) {
	L().Error(msg, tags...)
}

// Panicf 输出 panic 级别的日志.
func Panicf(format string, args ...any) {
	L().Panicf(format, args...)
}

// Panic 输出 panic 级别的日志.
func Panic(msg string, tags ...Field) {
	L().Panic(msg, tags...)
}

// Fatalf 输出 fatal 级别的日志.
func Fatalf(format string, args ...any) {
	L().Fatalf(format, args...)
}

// Fatal 输出 fatal 级别的日志.
func Fatal(msg string, tags ...Field) {
	L().Fatal(msg, tags...)
}

func With(fields ...Field) Logger {
	return L().With(fields...)
}

func AddCallerSkip(skip int) Logger {
	return L().AddCallerSkip(skip)
}
