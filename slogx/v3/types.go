package slogx

import (
	"context"
	"log/slog"
)

type Logger interface {
	Info(msg string, attrs ...Attr)
	Infof(format string, attrs ...any)
	InfoContext(ctx context.Context, msg string, attrs ...Attr)
	Error(msg string, attrs ...Attr)
	Errorf(format string, args ...any)
	ErrorContext(ctx context.Context, msg string, attrs ...Attr)
	Debug(msg string, attrs ...Attr)
	Debugf(format string, args ...any)
	DebugContext(ctx context.Context, msg string, attrs ...Attr)
	Warn(msg string, attrs ...Attr)
	Warnf(format string, args ...any)
	WarnContext(ctx context.Context, msg string, attrs ...Attr)
	Panic(msg string, attrs ...Attr)
	Panicf(format string, args ...any)
	PanicContext(ctx context.Context, msg string, attrs ...Attr)
	Fatal(msg string, attrs ...Attr)
	Fatalf(format string, args ...any)
	FatalContext(ctx context.Context, msg string, attrs ...Attr)
	Trace(msg string, attrs ...Attr)
	Tracef(format string, args ...any)
	TraceContext(ctx context.Context, msg string, attrs ...Attr)
	With(args ...any) *SlogLogger
	WithGroup(name string) *SlogLogger

	New(args ...any) *SlogLogger
	Enabled(ctx context.Context, level slog.Level) bool
}

type Attr = slog.Attr

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)
