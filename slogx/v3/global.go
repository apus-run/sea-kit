package slogx

import (
	"context"
	"log/slog"
	"sync/atomic"
)

var global atomic.Value

func init() {
	logger := NewLogger()

	global.Store(logger)
}

// SetLogger sets global logger appliance as logger in current process.
func SetLogger(logger Logger) {
	global.Store(logger)
	if lg, ok := logger.(*SlogLogger); ok {
		slog.SetDefault(lg.Logger)
	}
}

// GetLogger returns global logger appliance as logger in current process.
func GetLogger() Logger {
	return global.Load().(Logger)
}

// L 是 GetLogger 简写
func L() Logger {
	return global.Load().(Logger)
}

// Default 是 GetLogger
func Default() Logger {
	return global.Load().(Logger)
}

func Info(msg string, args ...Attr) {
	L().Info(msg, args...)
}

func Infof(format string, args ...any) {
	L().Infof(format, args...)
}

func InfoContext(ctx context.Context, msg string, args ...Attr) {
	L().InfoContext(ctx, msg, args...)
}

func Error(msg string, args ...Attr) {
	L().Error(msg, args...)
}

func Errorf(format string, args ...any) {
	L().Errorf(format, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...Attr) {
	L().ErrorContext(ctx, msg, args...)
}

func Debug(msg string, args ...Attr) {
	L().Debug(msg, args...)
}

func Debugf(format string, args ...any) {
	L().Debugf(format, args...)
}

func DebugContext(ctx context.Context, msg string, args ...Attr) {
	L().DebugContext(ctx, msg, args...)
}

func Warn(msg string, args ...Attr) {
	L().Warn(msg, args...)
}

func Warnf(format string, args ...any) {
	L().Warnf(format, args...)
}

func WarnContext(ctx context.Context, msg string, args ...Attr) {
	L().WarnContext(ctx, msg, args...)
}

func Fatal(msg string, args ...Attr) {
	L().Fatal(msg, args...)
}

func Fatalf(format string, args ...any) {
	L().Fatalf(format, args...)
}

func FatalContext(ctx context.Context, msg string, args ...Attr) {
	L().FatalContext(ctx, msg, args...)
}

func Panic(msg string, args ...Attr) {
	L().Panic(msg, args...)
}

func Panicf(format string, args ...any) {
	L().Panicf(format, args...)
}

func PanicContext(ctx context.Context, msg string, args ...Attr) {
	L().PanicContext(ctx, msg, args...)
}

func Trace(msg string, args ...Attr) {
	L().Trace(msg, args...)
}

func Tracef(format string, args ...any) {
	L().Tracef(format, args...)
}

func TraceContext(ctx context.Context, msg string, args ...Attr) {
	L().TraceContext(ctx, msg, args...)
}

func With(args ...any) *SlogLogger {
	return L().With(args...)
}
func New(args ...any) *SlogLogger {
	return With(args...)
}

func WithGroup(name string) *SlogLogger {
	return L().WithGroup(name)
}
