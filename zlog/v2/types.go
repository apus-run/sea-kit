package zlog

import (
	"context"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Panicf(format string, args ...any)

	Info(msg string, tags ...Field)
	Error(msg string, tags ...Field)
	Debug(msg string, tags ...Field)
	Warn(msg string, tags ...Field)
	Fatal(msg string, tags ...Field)
	Panic(msg string, tags ...Field)

	Slow(msg string, fields ...Field)
	Stack(msg string)
	Stat(msg string, fields ...Field)

	Print(args ...any)
	Printf(format string, args ...any)
	Println(args ...any)

	With(fields ...Field) Logger
	AddCallerSkip(skip int) Logger

	Close() error
	Sync() error

	WithContext(ctx context.Context, keyvals ...any) context.Context
}

type Field = zapcore.Field
