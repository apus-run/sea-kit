package requstid

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// ContextRequestIDKey context request id for context
	ContextRequestIDKey = "request_id"

	// HeaderXRequestIDKey http header request id key
	HeaderXRequestIDKey = "X-Request-ID"
)

// CtxKeyString for context.WithValue key type
type CtxKeyString string

// RequestIDKey "request_id"
var RequestIDKey = CtxKeyString(ContextRequestIDKey)

// RequestID is an interceptor that injects a 'X-Request-ID' into the context and request/response header of each request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.GetHeader(HeaderXRequestIDKey)

		// Create request id
		if requestID == "" {
			requestID := uuid.New().String()
			c.Request.Header.Set(HeaderXRequestIDKey, requestID)
			// Expose it for use in the application
			c.Set(ContextRequestIDKey, requestID)
		}

		// Set X-Request-ID header
		c.Writer.Header().Set(HeaderXRequestIDKey, requestID)

		c.Next()
	}
}

// GetCtxRequestID get request id from gin.Context
func GetCtxRequestID(c *gin.Context) string {
	if v, isExist := c.Get(ContextRequestIDKey); isExist {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}
	return ""
}

// CtxRequestID get request id from context.Context
func CtxRequestID(ctx context.Context) string {
	v := ctx.Value(ContextRequestIDKey)
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

// HeaderRequestID get request id from the header
func HeaderRequestID(c *gin.Context) string {
	return c.Request.Header.Get(HeaderXRequestIDKey)
}
