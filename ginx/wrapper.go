package ginx

import (
	"context"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/gin-gonic/gin"
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

// Context a wrapper of gin.Context
type Context struct {
	*gin.Context
}

// HandlerFunc defines the handler to wrap gin.Context
type HandlerFunc func(*Context)

// ProxyHandlerFunc represents the reverse proxy handler function
type ProxyHandlerFunc func(*Context) (*httputil.ReverseProxy, error)

// WrapContext returns a context wrapped by this file
func WrapContext(c *gin.Context) *Context {
	return &Context{
		Context: c,
	}
}

// Handle convert HandlerFunc to gin.HandlerFunc
func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := WrapContext(ginCtx)
		h(c)
	}
}

func ProxyHandle(fn ProxyHandlerFunc) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		c := WrapContext(ginCtx)
		proxy, err := fn(c)

		if err != nil {
			c.Data(http.StatusOK, "text/text", []byte(err.Error()))
			c.Abort()
		} else {
			proxy.ServeHTTP(c.Writer, c.Request)
		}
	}
}

type ginKey struct{}

// NewGinContext returns a new Context that carries gin.Context value.
func NewGinContext(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ginKey{}, c)
}

// FromGinContext returns the gin.Context value stored in ctx, if any.
func FromGinContext(ctx context.Context) (c *gin.Context, ok bool) {
	c, ok = ctx.Value(ginKey{}).(*gin.Context)
	return
}

// Result defines HTTP JSON response
type Result struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Data    any      `json:"data"`
	Details []string `json:"details,omitempty"`
}

// JSON returns JSON response
// e.x. {"code":<code>, "msg":<msg>, "data":<data>, "details":<details>}
func (ctx *Context) JSON(httpStatus int, resp Result) {
	ctx.Context.JSON(httpStatus, resp)
}

// JSONOK returns JSON response with successful business code and data
// e.x. {"code": 200, "msg":"成功", "data":<data>}
func (ctx *Context) JSONOK(msg string, data any) {
	j := new(Result)
	j.Code = CodeOK
	j.Msg = msg

	switch d := data.(type) {
	case error:
		j.Data = d.Error()
	case nil:
		j.Data = gin.H{}
	default:
		j.Data = data
	}

	ctx.Context.JSON(http.StatusOK, j)
}

// Success c.Success()
func (ctx *Context) Success(data ...any) {
	j := new(Result)
	j.Code = CodeOK
	j.Msg = "ok"

	if len(data) > 0 {
		j.Data = data[0]
	} else {
		j.Data = ""
	}

	ctx.Context.JSON(http.StatusOK, j)
}

// JSONE returns JSON response with failure business code ,msg and data
// e.x. {"code":<code>, "msg":<msg>, "data":<data>}
// c.JSONE(5, "系统错误", err)
func (ctx *Context) JSONE(code int, msg string, data any) {
	j := new(Result)
	j.Code = code
	j.Msg = msg

	switch d := data.(type) {
	case error:
		j.Data = d.Error()
	case nil:
		j.Data = gin.H{}
	default:
		j.Data = data
	}

	ctx.Context.JSON(http.StatusOK, j)
}

// Bind wraps gin context.Bind() with custom validator
func (ctx *Context) Bind(obj interface{}) (err error) {
	return validate(ctx.Context.Bind(obj))
}

// ShouldBind wraps gin context.ShouldBind() with custom validator
func (ctx *Context) ShouldBind(obj interface{}) (err error) {
	return validate(ctx.Context.ShouldBind(obj))
}

// NotFound 未找到相关路由
func (ctx *Context) NotFound() {
	ctx.String(http.StatusNotFound, "the route not found")
}

// GetClientLocale returns the client locale name
func (ctx *Context) GetClientLocale() string {
	value := ctx.GetHeader(AcceptLanguageHeaderName)

	return value
}

// GetClientTimezoneOffset returns the client timezone offset
func (ctx *Context) GetClientTimezoneOffset() (int16, error) {
	value := ctx.GetHeader(ClientTimezoneOffsetHeaderName)
	offset, err := strconv.Atoi(value)

	if err != nil {
		return 0, err
	}

	return int16(offset), nil
}

// SetRequestId sets the given request id to context
func (ctx *Context) SetRequestId(requestId string) {
	ctx.Set(requestIdFieldKey, requestId)
}

// GetRequestId returns the current request id
func (ctx *Context) GetRequestId() string {
	requestId, exists := ctx.Get(requestIdFieldKey)

	if !exists {
		return ""
	}

	return requestId.(string)
}
