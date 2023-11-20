package bbr

import (
	ratelimit "github.com/apus-run/sea-kit/ratelimit/bbr"
	"github.com/gin-gonic/gin"
)

type Builder struct {
	limiter ratelimit.Limiter
}

func NewBuilder(opts ...ratelimit.Option) *Builder {
	return &Builder{
		limiter: ratelimit.NewLimiter(opts...),
	}
}

func (b *Builder) Limiter(limiter ratelimit.Limiter) *Builder {
	b.limiter = limiter
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		done, err := b.limiter.Allow()
		if err != nil {
			// rejected
			c.AbortWithStatusJSON(429, gin.H{
				"code": 429,
				"msg":  "service unavailable due to rate limit exceeded",
			})
			return
		}

		// allowed
		done(ratelimit.DoneInfo{
			Err: c.Errors.Last(),
		})

		c.Next()
	}
}
