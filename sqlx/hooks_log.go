package sqlx

import (
	"context"
	"time"

	"github.com/qustavo/sqlhooks/v2"

	"github.com/apus-run/sea-kit/log"
)

var _ sqlhooks.Hooks = (*logHooks)(nil)
var _ sqlhooks.OnErrorer = (*logHooks)(nil)

// logHooks implement Hooks interface
type logHooks struct {
	log *log.Helper
}

func (h *logHooks) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	return context.WithValue(ctx, "started", time.Now()), nil
}

func (h *logHooks) After(ctx context.Context, query string, args ...any) (context.Context, error) {
	h.log.Infof("Query: `%s`, Args: `%q`. took: %s", query, args, time.Since(ctx.Value("started").(time.Time)))
	return ctx, nil
}

func (h *logHooks) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	h.log.Errorf("Error: %v, Query: `%s`, Args: `%q`, Took: %s", err, query, args, time.Since(ctx.Value("started").(time.Time)))
	return err
}
