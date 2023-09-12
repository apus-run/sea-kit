package zlog

import (
	"go.uber.org/zap"
)

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Panicf(format string, args ...any)

	Info(msg string, tags ...zap.Field)
	Error(msg string, tags ...zap.Field)
	Debug(msg string, tags ...zap.Field)
	Warn(msg string, tags ...zap.Field)
	Fatal(msg string, tags ...zap.Field)
	Panic(msg string, tags ...zap.Field)

	Close()
	Sync()
}
