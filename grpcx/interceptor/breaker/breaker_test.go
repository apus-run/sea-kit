package breaker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBuildUnaryServerInterceptor(t *testing.T) {
	ctx := context.Background()
	req := "test request"

	// 创建mock UnaryServerInfo 和 UnaryHandler
	info := &grpc.UnaryServerInfo{
		FullMethod: "/example.service/v1/SomeService/SomeMethod",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, status.New(codes.Internal, "Internal server error").Err()
	}

	builder := NewBreakerInterceptorBuilder()

	// 构建unary server拦截器
	interceptor := builder.BuildUnaryServerInterceptor()

	for i := 0; i < 110; i++ {
		resp, err := interceptor(ctx, req, info, handler)
		assert.Error(t, err)
		assert.Equal(t, nil, resp)

	}
	handler = func(ctx context.Context, req any) (any, error) {
		return nil, status.New(codes.PermissionDenied, "Permission").Err()
	}
	_, err := interceptor(ctx, req, info, handler)
	assert.Error(t, err)

}

func TestBuildStreamServerInterceptor(t *testing.T) {
	builder := NewBreakerInterceptorBuilder()
	interceptor := builder.BuildStreamServerInterceptor()

	handler := func(srv interface{}, stream grpc.ServerStream) error {
		return status.New(codes.Internal, "Internal server error").Err()
	}

	for i := 0; i < 110; i++ {
		err := interceptor(nil, nil, &grpc.StreamServerInfo{FullMethod: "/test"}, handler)
		assert.Error(t, err)
	}

	handler = func(srv interface{}, stream grpc.ServerStream) error {
		return status.New(codes.PermissionDenied, "Permission").Err()
	}
	err := interceptor(nil, nil, &grpc.StreamServerInfo{FullMethod: "/test"}, handler)
	assert.Error(t, err)
}

func TestBuildUnaryClientInterceptor(t *testing.T) {
	builder := NewBreakerInterceptorBuilder()
	builder.ValidCode(codes.PermissionDenied)
	interceptor := builder.BuildUnaryClientInterceptor()

	assert.NotNil(t, interceptor)

	ivoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return status.New(codes.Internal, "Internal server error").Err()
	}
	for i := 0; i < 110; i++ {
		err := interceptor(context.Background(), "/test", nil, nil, nil, ivoker)
		assert.Error(t, err)
	}

	ivoker = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return status.New(codes.PermissionDenied, "Permission").Err()
	}
	err := interceptor(context.Background(), "/test", nil, nil, nil, ivoker)
	assert.Error(t, err)
}

func TestBuildStreamClientInterceptor(t *testing.T) {
	builder := NewBreakerInterceptorBuilder()
	interceptor := builder.BuildStreamClientInterceptor()

	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, status.New(codes.Internal, "Internal server error").Err()
	}
	for i := 0; i < 110; i++ {
		_, err := interceptor(context.Background(), nil, nil, "/test", streamer)
		assert.Error(t, err)
	}

	streamer = func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, status.New(codes.PermissionDenied, "Permission").Err()
	}
	_, err := interceptor(context.Background(), nil, nil, "/test", streamer)
	assert.Error(t, err)
}
