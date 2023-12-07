package run_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/apus-run/sea-kit/concurrency/run"
)

func ExampleGroup_Add_basic() {
	var g run.Group
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			select {
			case <-time.After(time.Second):
				fmt.Printf("The first actor had its time elapsed\n")
				return nil
			case <-cancel:
				fmt.Printf("The first actor was canceled\n")
				return nil
			}
		}, func(err error) {
			fmt.Printf("The first actor was interrupted with: %v\n", err)
			close(cancel)
		})
	}
	{
		g.Add(func() error {
			fmt.Printf("The second actor is returning immediately\n")
			return errors.New("immediate teardown")
		}, func(err error) {
			// Note that this interrupt function is called, even though the
			// corresponding execute function has already returned.
			fmt.Printf("The second actor was interrupted with: %v\n", err)
		})
	}
	fmt.Printf("The group was terminated with: %v\n", g.Run())
	// Output:
	// The second actor is returning immediately
	// The first actor was interrupted with: immediate teardown
	// The second actor was interrupted with: immediate teardown
	// The first actor was canceled
	// The group was terminated with: immediate teardown
}

func ExampleGroup_Add_context() {
	ctx, cancel := context.WithCancel(context.Background())
	var g run.Group
	{
		ctx, cancel := context.WithCancel(ctx) // note: shadowed
		g.Add(func() error {
			return runUntilCanceled(ctx)
		}, func(error) {
			cancel()
		})
	}
	go cancel()
	fmt.Printf("The group was terminated with: %v\n", g.Run())
	// Output:
	// The group was terminated with: context canceled
}

func ExampleGroup_Add_listener() {
	var g run.Group
	{
		ln, _ := net.Listen("tcp", ":0")
		g.Add(func() error {
			defer fmt.Printf("http.Serve returned\n")
			return http.Serve(ln, http.NewServeMux())
		}, func(error) {
			ln.Close()
		})
	}
	{
		g.Add(func() error {
			return errors.New("immediate teardown")
		}, func(error) {
			//
		})
	}
	fmt.Printf("The group was terminated with: %v\n", g.Run())
	// Output:
	// http.Serve returned
	// The group was terminated with: immediate teardown
}

func runUntilCanceled(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

//// Server 启动微服务所需要的所有依赖
//type Server struct {
//	runGroup   *run.Group
//	grpcServer *grpc.Server
//}
//
//// NewServer 实例化 Server
//func NewServer(
//	runGroup *run.Group,
//	grpcServer *grpc.Server,
//) *Server {
//	return &Server{
//		runGroup:   runGroup,
//		grpcServer: grpcServer,
//	}
//}
//
//// RunServer 启动 http 以及 grpc 服务
//func (s *Server) RunServer() {
//	// 启动 grpc 服务
//	s.runGroup.Add(func() error {
//		l, err := net.Listen("tcp", ":8080")
//		if err != nil {
//			return err
//		}
//		log.Printf("starting gRPC server: %v", l.Addr().String())
//
//		return s.grpcServer.Serve(l)
//	}, func(err error) {
//		s.grpcServer.GracefulStop()
//		s.grpcServer.Stop()
//	})
//
//	// 监听退出信号
//	s.runGroup.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))
//
//	// 顺序启动服务
//	if err := s.runGroup.Run(); err != nil {
//		os.Exit(1)
//	}
//}
