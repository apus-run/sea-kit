package prometheus

import (
	"context"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/apus-run/sea-kit/grpcx/interceptor"
)

// Kind is the type of Interceptor
const Kind string = "Prometheus"

type InterceptorBuilder struct {
	Namespace string
	Subsystem string
	interceptor.Builder
}

func NewPrometheusInterceptorBuilder() *InterceptorBuilder {
	return &InterceptorBuilder{}
}

// Kind return the name of interceptor
func (b *InterceptorBuilder) Kind() string {
	return Kind
}

func (b *InterceptorBuilder) BuildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	// ServerHandleHistogram ...
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: b.Namespace,
			Subsystem: b.Subsystem,
			Name:      "server_handle_seconds",
			Objectives: map[float64]float64{
				0.5:   0.01,
				0.9:   0.01,
				0.95:  0.01,
				0.99:  0.001,
				0.999: 0.0001,
			},
		},
		[]string{"type", "service", "method", "peer", "code"},
	)
	prometheus.MustRegister(summary)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		defer func() {
			serviceName, method := b.splitMethodName(info.FullMethod)
			st, _ := status.FromError(err)
			code := "OK"
			if st != nil {
				code = st.Code().String()
			}
			summary.WithLabelValues("unary", serviceName, method,
				b.PeerName(ctx), code).Observe(float64(time.Since(start).Milliseconds()))
		}()
		resp, err = handler(ctx, req)
		return
	}
}

func (b *InterceptorBuilder) splitMethodName(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return "unknown", "unknown"
}
