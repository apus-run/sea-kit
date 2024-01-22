package server

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/apus-run/sea-kit/grpcx/client"
	"github.com/apus-run/sea-kit/grpcx/interceptor/recovery"
	pb "github.com/apus-run/sea-kit/grpcx/testdata/helloworld"
	"google.golang.org/grpc"
)

// service is used to implement helloworld.GreeterServer.
type service struct {
	pb.UnimplementedGreeterServer
}

func (s *service) SayHelloStream(streamServer pb.Greeter_SayHelloStreamServer) error {
	var cnt uint
	for {
		in, err := streamServer.Recv()
		if err != nil {
			return err
		}
		if in.Name == "error" {
			panic(fmt.Sprintf("invalid argument %s", in.Name))
		}
		if in.Name == "panic" {
			panic("server panic")
		}
		err = streamServer.Send(&pb.HelloReply{
			Message: fmt.Sprintf("hello %s", in.Name),
		})
		if err != nil {
			return err
		}
		cnt++
		if cnt > 1 {
			return nil
		}
	}
}

// SayHello implements helloworld.GreeterServer
func (s *service) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Name == "error" {
		panic(fmt.Sprintf("invalid argument %s", in.Name))
	}
	if in.Name == "panic" {
		panic("server panic")
	}
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

type testKey struct{}

func TestServer(t *testing.T) {
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	ctx := context.Background()
	ctx = context.WithValue(ctx, testKey{}, "test")

	interceptor := recovery.NewRecoveryInterceptorBuilder()
	srv := NewServer(
		WithAddr(":8090"),
		WithUnaryInterceptor(
			func(ctx context.Context, req interface{},
				info *grpc.UnaryServerInfo,
				handler grpc.UnaryHandler) (resp interface{}, err error) {
				return handler(ctx, req)
			},
			interceptor.BuildUnaryServerInterceptor(),
		),
		WithStreamInterceptor(
			interceptor.BuildStreamServerInterceptor(),
		),
		WithGrpcOptions(grpc.InitialConnWindowSize(0)),
	)
	pb.RegisterGreeterServer(srv, &service{})

	go func() {
		// start server
		if err := srv.Start(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	testClient(t, ctx, srv)
	srv.Stop()
}

func testClient(t *testing.T, ctx context.Context, srv *Options) {
	u, err := srv.Endpoint()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)

	// new a gRPC client
	conn, err := client.NewClient(
		ctx,
		client.WithAddr(":8090"),
		client.WithUnaryInterceptor(func(
			ctx context.Context,
			method string, req,
			reply interface{},
			cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker,
			opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	defer conn.Close()
	if err != nil {
		t.Fatal(err)
	}
	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "gaea"})
	t.Log(err)
	if err != nil {
		t.Errorf("failed to call: %v", err)
	}
	if !reflect.DeepEqual(reply.Message, "Hello gaea") {
		t.Errorf("expect %s, got %s", "Hello gaea", reply.Message)
	}

	streamCli, err := client.SayHelloStream(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = streamCli.CloseSend()
	}()
	err = streamCli.Send(&pb.HelloRequest{Name: "cc"})
	if err != nil {
		t.Error(err)
		return
	}
	reply, err = streamCli.Recv()
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(reply.Message, "hello cc") {
		t.Errorf("expect %s, got %s", "hello cc", reply.Message)
	}
}

func TestNewServer(t *testing.T) {
	ctx := context.Background()
	srv := NewServer(
		WithAddr(":8090"),
	)

	// Attach the Greeter service to the server
	pb.RegisterGreeterServer(srv, &service{})
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8090")
	go func() {
		// start server
		if err := srv.Start(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	// new a gRPC client
	conn, err := client.NewClient(
		ctx,
		client.WithAddr(":8090"),
		client.WithUnaryInterceptor(func(
			ctx context.Context,
			method string, req,
			reply interface{},
			cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker,
			opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	defer conn.Close()

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "gaea"})
	t.Log(err)
	if err != nil {
		t.Errorf("failed to call: %v", err)
	}
	if !reflect.DeepEqual(reply.Message, "Hello gaea") {
		t.Errorf("expect %s, got %s", "Hello gaea", reply.Message)
	}

	streamCli, err := client.SayHelloStream(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = streamCli.CloseSend()
	}()
	err = streamCli.Send(&pb.HelloRequest{Name: "cc"})
	if err != nil {
		t.Error(err)
		return
	}
	reply, err = streamCli.Recv()
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(reply.Message, "hello cc") {
		t.Errorf("expect %s, got %s", "hello cc", reply.Message)
	}

	t.Log("输出:", reply.Message)

	srv.Stop()
}
