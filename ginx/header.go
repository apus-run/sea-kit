package ginx

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NoCache 控制客户端不要使用缓存
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, max-age=0, must-revalidate")
		c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		c.Next()
	}
}

// Secure 添加安全控制和资源访问
func Secure() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000")
		}
		c.Next()
	}
}
