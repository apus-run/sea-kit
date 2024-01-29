package log

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/apus-run/sea-kit/zlog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestBuildUnaryServerInterceptor(t *testing.T) {
	ctx := context.Background()
	req := "test request"

	// 创建mock UnaryServerInfo 和 UnaryHandler
	info := &grpc.UnaryServerInfo{
		FullMethod: "/example.service/v1/SomeService/SomeMethod",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return "test response", nil
	}

	// 创建LoggerInterceptorBuilder实例
	builder := NewLoggerInterceptorBuilder(zlog.L())

	// 构建unary server拦截器
	interceptor := builder.BuildUnaryServerInterceptor()

	// 模拟正常调用情况
	resp, err := interceptor(ctx, req, info, handler)
	assert.NoError(t, err)
	assert.Equal(t, "test response", resp)

	// 检查是否跳过了健康检查方法
	healthCheckInfo := &grpc.UnaryServerInfo{FullMethod: "/grpc.health.v1.Health/Check"}
	resp, err = interceptor(ctx, req, healthCheckInfo, handler)
	assert.NoError(t, err)
	assert.Equal(t, "test response", resp)

	// 模拟恢复panic的情况
	recoverInfo := &grpc.UnaryServerInfo{FullMethod: "/example.service/v1/SomeService/PanicMethod"}
	panicHandler := func(ctx context.Context, req any) (any, error) {
		panic(errors.New("forced panic"))
	}
	_, err = interceptor(ctx, req, recoverInfo, panicHandler)
	assert.Error(t, err)
	t.Logf("err: %v", err.Error())
}

func TestBuildUnaryClientInterceptor(t *testing.T) {
	// 创建LoggerInterceptorBuilder实例
	builder := NewLoggerInterceptorBuilder(zlog.L())
	builder.IgnoreMethods("/test")
	// 构建unary client拦截器
	interceptor := builder.BuildUnaryClientInterceptor()

	assert.NotNil(t, interceptor)

	ivoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return status.New(codes.Internal, "Internal server error").Err()
	}

	err := interceptor(context.Background(), "/test", nil, nil, nil, ivoker)
	assert.Error(t, err)
	t.Logf("err: %v", err.Error())
}
