package sqlx

import (
	"context"
	"time"

	"github.com/qustavo/sqlhooks/v2"
)

var _ sqlhooks.Hooks = (*metricHooks)(nil)
var _ sqlhooks.OnErrorer = (*metricHooks)(nil)

// metricHooks implement Hooks interface
type metricHooks struct {
}

func (h *metricHooks) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	return context.WithValue(ctx, "metric started", time.Now()), nil
}

func (h *metricHooks) After(ctx context.Context, query string, args ...any) (context.Context, error) {

	return ctx, nil
}

func (h *metricHooks) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	return err
}
