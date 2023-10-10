package accesslog

import (
	"bytes"
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
)

// Builder 记录HTTP请求/响应细节
type Builder struct {
	allowReqBody  *atomic.Bool
	allowRespBody *atomic.Bool

	// http response body 的 max length; request URL 的 max length.
	maxLength *atomic.Int64

	// 忽略指定路由的日志打印
	ignoreRoutes map[string]struct{}

	logFunc func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *Builder {
	return &Builder{
		// 默认不打印
		allowReqBody:  atomic.NewBool(false),
		allowRespBody: atomic.NewBool(false),

		maxLength: atomic.NewInt64(1024), // 1 MiB

		ignoreRoutes: map[string]struct{}{
			"/ping":   {},
			"/pong":   {},
			"/health": {},
		},

		logFunc: fn,
	}
}

func (b *Builder) AllowReqBody() *Builder {
	b.allowReqBody.Store(true)
	return b
}

func (b *Builder) AllowRespBody() *Builder {
	b.allowRespBody.Store(true)
	return b
}

func (b *Builder) MaxLength(maxLength int64) *Builder {
	b.maxLength.Store(maxLength)
	return b
}

func (b *Builder) IgnoreRoutes(routes ...string) *Builder {
	for _, route := range routes {
		b.ignoreRoutes[route] = struct{}{}
	}

	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	pid := strconv.Itoa(os.Getpid())
	return func(c *gin.Context) {
		start := time.Now()
		maxLength := b.maxLength.Load()
		allowReqBody := b.allowReqBody.Load()
		allowRespBody := b.allowRespBody.Load()

		// ignore printing of the specified route
		if _, ok := b.ignoreRoutes[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		host := c.Request.Host
		split := strings.Split(host, ":")

		// URL 有可能会很长, 保护起来
		url := c.Request.URL
		urlStr := url.String()
		urlLen := int64(len(urlStr))
		if urlLen >= maxLength {
			urlStr = urlStr[:maxLength]
		}
		accessLog := &AccessLog{
			PID:      pid,
			Referer:  c.Request.Header.Get("Referer"),
			Protocol: url.Scheme,
			Port:     split[1],
			IP:       split[0],
			IPs:      c.Request.Header.Get("X-Forwarded-For"),
			Host:     host,
			URL:      urlStr,
			UA:       c.Request.Header.Get("User-Agent"),

			Method: c.Request.Method,
			Path:   url.Path,
		}

		if allowReqBody && c.Request.Body != nil {
			// 可以直接忽略 error，不影响程序运行
			// GetRawData 实现了 io.ReadAll(c.Request.Body)
			body, _ := c.GetRawData()
			// Request.Body 是一个 Stream（流）对象，所以是只能读取一次的
			// 因此读完之后要放回去，不然后续步骤是读不到的
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// 防止body内容过大, 保护起来
			if int64(len(body)) >= maxLength {
				body = body[:maxLength]
			}
			//注意资源的消耗
			accessLog.ReqBody = string(body)
		}

		if allowRespBody {
			c.Writer = responseWriter{
				al:             accessLog,
				ResponseWriter: c.Writer,
				maxLength:      maxLength,
			}
		}

		defer func() {
			duration := time.Since(start)
			accessLog.Duration = duration.String()
			b.logFunc(c, accessLog)
		}()

		c.Next()
	}
}
