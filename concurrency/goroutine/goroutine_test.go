package goroutine

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	pkggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/apus-run/sea-kit/concurrency/goroutine/proto/helloworld"
)

// Server is used to implement helloworld.GreeterServer.
type Server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func NewGrpcServer() (*pkggrpc.Server, error) {
	s := pkggrpc.NewServer()

	// 这里进行服务注册
	pb.RegisterGreeterServer(s, &Server{})
	reflection.Register(s)

	return s, nil
}

// 启动AppServer, 这个函数会将当前goroutine阻塞
func startAppServe(ctx context.Context, server *grpc.Server, lis net.Listener) error {
	// 这个goroutine是启动服务的goroutine
	SafeGo(ctx, func() {
		if err := server.Serve(lis); err != nil {
			fmt.Errorf("grpc serve error %v", map[string]interface{}{
				"error": err.Error(),
			})
		}
	})

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	server.GracefulStop()

	return nil
}

func TestStartServer(t *testing.T) {
	server, err := NewGrpcServer()
	if err != nil {
		return
	}

	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx := context.Background()
	if err := startAppServe(ctx, server, lis); err != nil {
		fmt.Println(err)
	}
}

// TestMultiHandlerDemo goroutine 的使用示例
func TestMultiHandlerDemoGoroutine(t *testing.T) {
	type User struct {
		ID   uint64
		Name string
	}
	// 初始化一个orm.DB
	db := gorm.DB{}
	c := context.Background()

	err := SafeGoAndWait(c, func() error {
		// 查询一条数据
		queryUser := &User{ID: 1}

		err := db.First(queryUser).Error

		fmt.Printf("query user1 %v", map[string]interface{}{
			"err":  err,
			"name": queryUser.Name,
		})
		return err
	}, func() error {
		// 查询一条数据
		queryUser := &User{ID: 2}

		err := db.First(queryUser).Error
		fmt.Printf("query user2 %v", map[string]interface{}{
			"err":  err,
			"name": queryUser.Name,
		})
		return err
	})

	if err != nil {
		fmt.Errorf("error: %v", err)
	}
}

func TestSafeGo(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	SafeGo(ctx, func() {
		time.Sleep(1 * time.Second)
		return
	})
	t.Log("safe go main start")
	time.Sleep(2 * time.Second)
	t.Log("safe go main end")

	SafeGo(ctx, func() {
		time.Sleep(1 * time.Second)
		panic("safe go test panic")
	})
	t.Log("safe go2 main start")
	time.Sleep(2 * time.Second)
	t.Log("safe go2 main end")

}

func TestSafeGoAndWait(t *testing.T) {
	errStr := "safe go test error"
	t.Log("safe go and wait start", time.Now().String())
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	err := SafeGoAndWait(ctx, func() error {
		time.Sleep(1 * time.Second)
		return errors.New(errStr)
	}, func() error {
		time.Sleep(2 * time.Second)
		return nil
	}, func() error {
		time.Sleep(3 * time.Second)
		return nil
	})
	t.Log("safe go and wait end", time.Now().String())

	if err == nil {
		t.Error("err not be nil")
	} else if err.Error() != errStr {
		t.Error("err content not same")
	}

	// panic error
	err = SafeGoAndWait(ctx, func() error {
		time.Sleep(1 * time.Second)
		return errors.New(errStr)
	}, func() error {
		time.Sleep(2 * time.Second)
		panic("test2")
	}, func() error {
		time.Sleep(3 * time.Second)
		return nil
	})
	if err == nil {
		t.Error("err not be nil")
	} else if err.Error() != errStr {
		t.Error("err content not same")
	}
}
