package hooks

import (
	"context"
	"time"

	"github.com/qustavo/sqlhooks/v2"
)

var _ sqlhooks.Hooks = (*MetricHooks)(nil)
var _ sqlhooks.OnErrorer = (*MetricHooks)(nil)

// MetricHooks implement Hooks interface
type MetricHooks struct {
}

func (h *MetricHooks) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	return context.WithValue(ctx, "metric started", time.Now()), nil
}

func (h *MetricHooks) After(ctx context.Context, query string, args ...any) (context.Context, error) {

	return ctx, nil
}

func (h *MetricHooks) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	return err
}
