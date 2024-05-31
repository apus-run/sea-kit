package hooks

import (
	"context"
	"time"

	"github.com/qustavo/sqlhooks/v2"

	"github.com/apus-run/sea-kit/log"
)

var _ sqlhooks.Hooks = (*LogHooks)(nil)
var _ sqlhooks.OnErrorer = (*LogHooks)(nil)

// logHooks implement Hooks interface
type LogHooks struct {
	Log *log.Helper
}

func (h *LogHooks) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	return context.WithValue(ctx, "started", time.Now()), nil
}

func (h *LogHooks) After(ctx context.Context, query string, args ...any) (context.Context, error) {
	h.Log.Infof("Query: `%s`, Args: `%q`. took: %s", query, args, time.Since(ctx.Value("started").(time.Time)))
	return ctx, nil
}

func (h *LogHooks) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	h.Log.Errorf("Error: %v, Query: `%s`, Args: `%q`, Took: %s", err, query, args, time.Since(ctx.Value("started").(time.Time)))
	return err
}
