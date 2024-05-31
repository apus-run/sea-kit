package hooks

import (
	"context"
	"time"

	"github.com/qustavo/sqlhooks/v2"
)

var _ sqlhooks.Hooks = (*TracingHooks)(nil)
var _ sqlhooks.OnErrorer = (*TracingHooks)(nil)

// TracingHooks implement Hooks interface
type TracingHooks struct{}

func (h *TracingHooks) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	return context.WithValue(ctx, "tracing started", time.Now()), nil
}

func (h *TracingHooks) After(ctx context.Context, query string, args ...any) (context.Context, error) {

	return ctx, nil
}

func (h *TracingHooks) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	return err
}
