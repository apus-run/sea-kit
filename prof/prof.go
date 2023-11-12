// Package prof is used for gin profiling.
package prof

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/felixge/fgprof"
	"github.com/gin-gonic/gin"
)

var defaultPrefix = "/debug/pprof"

// Option set defaultPrefix func
type Option func(o *options)

type options struct {
	prefix           string
	enableIOWaitTime bool
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithPrefix set route defaultPrefix
func WithPrefix(prefix string) Option {
	return func(o *options) {
		if prefix == "" {
			return
		}
		o.prefix = prefix
	}
}

// WithIOWaitTime enable IO wait time
func WithIOWaitTime() Option {
	return func(o *options) {
		o.enableIOWaitTime = true
	}
}

// Register pprof for gin router
func Register(r *gin.Engine, opts ...Option) {
	o := &options{prefix: defaultPrefix}
	o.apply(opts...)

	group := r.Group(o.prefix)

	group.GET("/", gin.WrapF(pprof.Index))
	group.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	group.GET("/profile", gin.WrapF(pprof.Profile))
	group.POST("/symbol", gin.WrapF(pprof.Symbol))
	group.GET("/symbol", gin.WrapF(pprof.Symbol))
	group.GET("/trace", gin.WrapF(pprof.Trace))
	group.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
	group.GET("/block", gin.WrapH(pprof.Handler("block")))
	group.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
	group.GET("/heap", gin.WrapH(pprof.Handler("heap")))
	group.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
	group.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))

	if o.enableIOWaitTime {
		// Similar to /profile, add IO wait time,  https://github.com/felixge/fgprof
		group.GET("/profile-io", gin.WrapH(fgprof.Handler()))
	}
}

// New 创建一个http ServeMux实例
func New() *http.ServeMux {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/debug/pprof/", pprof.Index)
	httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	httpMux.HandleFunc("/check", Check)

	return httpMux
}

// Run 运行PProf性能监控服务
// 性能监控的端口port只能在内网访问
// 一般在程序启动init或main函数中执行
// httpMux 表示*http.ServeMux
// port表示pprof运行的端口
func Run(httpMux *http.ServeMux, port int) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("PProf exec recover: ", err)
			}
		}()

		log.Println("server PProf run on: ", port)

		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), httpMux); err != nil {
			log.Println("PProf listen error: ", err)
		}

	}()

}

// Check PProf心跳检测
func Check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"alive": true}`))
}
