package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
)

const (
	// CodeOK means a successful response
	CodeOK = 0
	// CodeErr means a failure response
	CodeErr = 1
)

const requestIdFieldKey = "REQUEST_ID"

// AcceptLanguageHeaderName represents the header name of accept language
const AcceptLanguageHeaderName = "Accept-Language"

// ClientTimezoneOffsetHeaderName represents the header name of client timezone offset
const ClientTimezoneOffsetHeaderName = "X-Timezone-Offset"

type Handler interface {
	PrivateRoutes(server *gin.Engine)
	PublicRoutes(server *gin.Engine)
}

// Context a wrapper of gin.Context
type Context struct {
	*gin.Context
}

// HandlerFunc defines the handler to wrap gin.Context
type HandlerFunc func(*Context)

// ProxyHandlerFunc represents the reverse proxy handler function
type ProxyHandlerFunc func(*Context) (*httputil.ReverseProxy, error)

// Result defines HTTP JSON response
type Result struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Data    any      `json:"data"`
	Details []string `json:"details,omitempty"`
}
