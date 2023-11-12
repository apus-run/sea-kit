package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/apus-run/sea-kit/log"
	"github.com/apus-run/sea-kit/utils"
)

const (
	ResultSuccess = "success"
	ResultError   = "error"
)

// web程序的性能监控，如果是job/rpc服务就不需要这两行
// prometheus.MustRegister(WebRequestTotal)
// prometheus.MustRegister(WebRequestDuration)
var (
	// WebRequestTotal 初始化 web_request_total， counter类型指标， 表示接收http请求总次数
	// 设置两个标签 请求方法和 路径 对请求总次数在两个
	WebRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "web_request_total",
			Help: "Number of hello requests in total",
		},
		[]string{"method", "endpoint"},
	)

	// WebRequestDuration web_request_duration_seconds，
	// Histogram类型指标，bucket代表duration的分布区间
	WebRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "web_request_duration_seconds",
			Help:    "web request duration distribution",
			Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1},
		},
		[]string{"method", "endpoint"},
	)

	// CpuTemp cpu情况
	CpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU",
	})

	HdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors",
		},
		[]string{"device"},
	)
)

// MonitorHandlerFunc 对于http原始的处理器函数，包装 handler function,不侵入业务逻辑
// 可以对单个接口做metrics监控
func MonitorHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h(w, r)

		// counter类型 metric 的记录方式
		WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		// Histogram类型 metric的记录方式
		WebRequestDuration.With(prometheus.Labels{
			"method": r.Method, "endpoint": r.URL.Path,
		}).Observe(time.Since(start).Seconds())
	}
}

// MonitorHandler 性能监控处理器
// 可以作为中间件对接口进行打点监控
func MonitorHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		// counter类型 metric 的记录方式
		WebRequestTotal.With(prometheus.Labels{"method": r.Method, "endpoint": r.URL.Path}).Inc()
		// Histogram类型 metric 的记录方式
		WebRequestDuration.With(prometheus.Labels{
			"method": r.Method, "endpoint": r.URL.Path,
		}).Observe(time.Since(start).Seconds())
	})
}

var (
	SQLTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kine_sql_total",
		Help: "Total number of SQL operations",
	}, []string{"error_code"})

	SQLTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "kine_sql_time_seconds",
		Help: "Length of time per SQL operation",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.15, 0.2, 0.25, 0.3, 0.35, 0.4, 0.45, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0,
			1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5, 6, 7, 8, 9, 10, 15, 20, 25, 30},
	}, []string{"error_code"})

	CompactTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kine_compact_total",
		Help: "Total number of compactions",
	}, []string{"result"})
)

var (
	// SlowSQLThreshold is a duration which SQL executed longer than will be logged.
	// This can be directly modified to override the default value when kine is used as a library.
	SlowSQLThreshold = time.Second
)

func ObserveSQL(start time.Time, errCode string, sql utils.Stripped, args ...interface{}) {
	SQLTotal.WithLabelValues(errCode).Inc()
	duration := time.Since(start)
	SQLTime.WithLabelValues(errCode).Observe(duration.Seconds())
	if duration >= SlowSQLThreshold {
		log.Infof("Slow SQL (started: %v) (total time: %v): %s : %v", start, duration, sql, args)
	}
}
