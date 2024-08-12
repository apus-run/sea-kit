package timeout

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Builder struct {
	timeout time.Duration
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) SetTimeout(timeout time.Duration) *Builder {
	b.timeout = timeout
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), b.timeout)

		defer func() {
			//cancel to clear resources after finished
			cancel()

			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {

				// 记录操作日志
				log.Panicln(c, "server timeout", nil)

				// write response and abort the request
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
					"code":    504,
					"message": http.StatusText(http.StatusGatewayTimeout),
				})

			}
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// cost.NewBuilder().Build()
