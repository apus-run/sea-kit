package slogx

import (
	"context"
	"log/slog"
)

type ContextLogKey struct{}

func NewContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ContextLogKey{}, l)
}

func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	if _, ok := ctx.Value(ContextLogKey{}).(*slog.Logger); ok {
		return ctx
	}
	return context.WithValue(ctx, ContextLogKey{}, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ContextLogKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

//func log(ctx context.Context, level slog.Level, msg string, args ...any) {
//	FromContext(ctx).Log(ctx, level, msg, args...)
//}
