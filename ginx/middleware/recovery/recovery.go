package recovery

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// log.Println("error: ", err)
				// log.Println("exec panic error", safemap[string]interface{}{
				// 	"trace_error": string(debug.Stack()),
				// })

				ctx := c.Request.Context()
				log.Println("exec panic error", map[string]interface{}{
					"module":      "web",
					"trace_error": string(debug.Stack()),
				})

				// broker pipe
				if isBrokenPipe(ctx, err) {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				// services error
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
					"code":    http.StatusInternalServerError,
					"message": "server inner error",
				})
				return
			}
		}()

		c.Next()
	}
}

func isBrokenPipe(ctx context.Context, err interface{}) bool {
	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	var brokenPipe bool
	if ne, ok := err.(*net.OpError); ok {
		if se, exist := ne.Err.(*os.SyscallError); exist {
			errMsg := strings.ToLower(se.Error())
			// logger error
			log.Println(ctx, "os syscall error", map[string]interface{}{
				"trace_error": errMsg,
			})

			if strings.Contains(errMsg, "broken pipe") ||
				strings.Contains(errMsg, "reset by peer") ||
				strings.Contains(errMsg, "request headers: small read buffer") ||
				strings.Contains(errMsg, "unexpected EOF") ||
				strings.Contains(errMsg, "i/o timeout") {
				brokenPipe = true
			}
		}
	}

	return brokenPipe
}
