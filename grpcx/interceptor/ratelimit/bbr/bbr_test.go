package bbr

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	ratelimit "github.com/apus-run/sea-kit/ratelimit/bbr"
	"github.com/apus-run/sea-kit/zlog"
)

func TestBuildUnaryServerInterceptor(t *testing.T) {
	limiter := ratelimit.NewLimiter(
		ratelimit.WithWindow(5*time.Second),
		ratelimit.WithBucket(50),
		ratelimit.WithCPUThreshold(100),
	)
	builder := NewRatelimitInterceptorBuilder(limiter, zlog.L())
	interceptor := builder.BuildUnaryServerInterceptor()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}
	_, err := interceptor(nil, nil, nil, handler)
	assert.NoError(t, err)
}

func TestBuildStreamServerInterceptor(t *testing.T) {
	limiter := ratelimit.NewLimiter(
		ratelimit.WithWindow(5*time.Second),
		ratelimit.WithBucket(50),
		ratelimit.WithCPUThreshold(100),
	)
	builder := NewRatelimitInterceptorBuilder(limiter, zlog.L())
	interceptor := builder.BuildStreamServerInterceptor()

	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}
	err := interceptor(nil, nil, nil, handler)
	assert.NoError(t, err)
}
