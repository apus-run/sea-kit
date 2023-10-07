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
)

// Builder 记录HTTP请求/响应细节
type Builder struct {
	allowReqBody  bool
	allowRespBody bool

	// http response body 的 max length; request URL 的 max length.
	maxLength int

	// 忽略指定路由的日志打印
	ignoreRoutes map[string]struct{}

	logFunc func(ctx context.Context, al AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al AccessLog)) *Builder {
	return &Builder{
		// 默认不打印
		allowReqBody:  false,
		allowRespBody: false,

		maxLength: 1 << 20, // 1 MiB

		ignoreRoutes: map[string]struct{}{
			"/ping":   {},
			"/pong":   {},
			"/health": {},
		},

		logFunc: fn,
	}
}

func (b *Builder) AllowReqBody() *Builder {
	b.allowReqBody = true
	return b
}

func (b *Builder) AllowRespBody() *Builder {
	b.allowRespBody = true
	return b
}

func (b *Builder) SetMaxLength(maxLength int) *Builder {
	b.maxLength = maxLength
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
		if len(urlStr) > b.maxLength {
			urlStr = urlStr[:b.maxLength]
		}
		al := AccessLog{
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

		if b.allowReqBody && c.Request.Body != nil {
			// 可以直接忽略 error，不影响程序运行
			// GetRawData 实现了 io.ReadAll(c.Request.Body)
			body, _ := c.GetRawData()
			// Request.Body 是一个 Stream（流）对象，所以是只能读取一次的
			// 因此读完之后要放回去，不然后续步骤是读不到的
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// 防止body内容过大, 保护起来
			if len(body) > b.maxLength {
				body = body[:b.maxLength]
			}
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			c.Writer = responseWriter{
				al:             &al,
				ResponseWriter: c.Writer,
			}
		}

		defer func() {
			duration := time.Since(start)
			al.Duration = duration.String()
			b.logFunc(c, al)
		}()

		c.Next()
	}
}
