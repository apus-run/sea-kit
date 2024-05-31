package sqlx

import (
	"context"
	"time"

	"github.com/qustavo/sqlhooks/v2"
)

var _ sqlhooks.Hooks = (*tracingHooks)(nil)
var _ sqlhooks.OnErrorer = (*tracingHooks)(nil)

// tracingHooks implement Hooks interface
type tracingHooks struct{}

func (h *tracingHooks) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	return context.WithValue(ctx, "tracing started", time.Now()), nil
}

func (h *tracingHooks) After(ctx context.Context, query string, args ...any) (context.Context, error) {

	return ctx, nil
}

func (h *tracingHooks) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	return err
}
