package redis

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ratelimit "github.com/apus-run/sea-kit/ratelimit/redis"
	"github.com/apus-run/sea-kit/zlog"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestBuildUnaryServerInterceptor(t *testing.T) {
	ctx := context.Background()
	req := "test request"
	r := ratelimit.NewRedisSlidingWindowLimiter(
		initRedis(),
		500*time.Millisecond,
		5,
	)
	assert.NotNil(t, r)

	builder := NewRatelimitInterceptorBuilder(r, "foo", zlog.L())
	interceptor := builder.BuildUnaryServerInterceptor()

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, status.New(codes.Internal, "Internal server error").Err()
	}

	for i := 0; i < 6; i++ {
		resp, err := interceptor(ctx, req, info, handler)
		assert.Error(t, err)
		assert.Equal(t, nil, resp)
	}
}

func TestBuildUnaryServerInterceptorService(t *testing.T) {
	ctx := context.Background()
	req := "test request"
	r := ratelimit.NewRedisSlidingWindowLimiter(
		initRedis(),
		500*time.Millisecond,
		1,
	)
	assert.NotNil(t, r)

	builder := NewRatelimitInterceptorBuilder(r, "foo", zlog.L())
	interceptor := builder.BuildUnaryServerInterceptorService("/test")

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, status.New(codes.Internal, "Internal server error").Err()
	}

	resp, err := interceptor(ctx, req, info, handler)
	assert.Error(t, err)
	assert.Equal(t, nil, resp)
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		DB:       1,
		Password: "123456",
	})
	return redisClient
}
