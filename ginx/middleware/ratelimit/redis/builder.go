package redis

import (
	"log"
	"net/http"
	"strings"

	ratelimit "github.com/apus-run/sea-kit/ratelimit/redis"
	"github.com/gin-gonic/gin"
)

type Builder struct {
	limiter ratelimit.Limiter

	genKeyFn func(ctx *gin.Context) string
}

func NewBuilder(limiter ratelimit.Limiter) *Builder {
	return &Builder{
		limiter: limiter,
		genKeyFn: func(ctx *gin.Context) string {
			var b strings.Builder
			b.WriteString("ip-limiter")
			b.WriteString(":")
			b.WriteString(ctx.ClientIP())
			return b.String()
		},
	}
}

func (b *Builder) SetKeyGenFunc(fn func(*gin.Context) string) *Builder {
	b.genKeyFn = fn
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	return b.limiter.Limit(ctx, b.genKeyFn(ctx))
}
