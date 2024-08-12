package secure

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"

	"log"
)

type Builder struct {
	*secure.Options
}

func (b *Builder) WithSecureOptions(opts *secure.Options) *Builder {
	b.Options = opts
	return b
}

func NewBuilder() *Builder {
	return &Builder{
		Options: &secure.Options{
			SSLRedirect: true,
			SSLHost:     "127.0.0.1:443",
		},
	}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(*b.Options)
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			log.Println(err)
			return
		}
		c.Next()
	}
}
