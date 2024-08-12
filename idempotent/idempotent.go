package idempotent

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// redis lua script(read => delete => get delete flag).
const (
	lua string = `
local current = redis.call('GET', KEYS[1])
if current == false then
    return '-1';
end
local del = redis.call('DEL', KEYS[1])
if del == 1 then
     return '1';
else
     return '0';
end
`
)

type Idempotent struct {
	ops Options
}

func New(options ...func(*Options)) *Idempotent {
	ops := getOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	return &Idempotent{ops: *ops}
}

func (i *Idempotent) Token(ctx context.Context) string {
	if i.ops.redis == nil {
		slog.Default().WarnContext(ctx, "please enable redis, otherwise the idempotent is invalid")
		return ""
	}

	token := uuid.NewString()
	i.ops.redis.Set(ctx, fmt.Sprintf("%s_%s", i.ops.prefix, token), true, time.Duration(i.ops.expire)*time.Minute)
	return token
}

func (i *Idempotent) Check(ctx context.Context, token string) bool {
	if i.ops.redis == nil {
		slog.Default().WarnContext(ctx, "please enable redis, otherwise the idempotent is invalid")
		return true
	}

	res, err := i.ops.redis.Eval(ctx, lua, []string{fmt.Sprintf("%s_%s", i.ops.prefix, token)}).Result()
	if err != nil || res != "1" {
		return false
	}

	return true
}
