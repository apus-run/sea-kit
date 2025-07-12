package slogx

import (
	"context"
)

type ContextLogKey struct{}

func NewContext(ctx context.Context, l *SlogLogger) context.Context {
	return context.WithValue(ctx, ContextLogKey{}, l)
}

func WithContext(ctx context.Context, l *SlogLogger) context.Context {
	if _, ok := ctx.Value(ContextLogKey{}).(*SlogLogger); ok {
		return ctx
	}
	return context.WithValue(ctx, ContextLogKey{}, l)
}

func FromContext(ctx context.Context) *SlogLogger {
	if l, ok := ctx.Value(ContextLogKey{}).(*SlogLogger); ok {
		return l
	}
	return nil
}

// C represents for `FromContext` with empty keyvals.
// slog.C(ctx).Info("Set function called")
func C(ctx context.Context) *SlogLogger {
	return FromContext(ctx)
}
