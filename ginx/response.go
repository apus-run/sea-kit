package ginx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

// Response 代表返回方法
type Response interface {
	// Json输出
	Json(obj interface{}) Response

	// Jsonp输出
	Jsonp(obj interface{}) Response

	//xml输出
	Xml(obj interface{}) Response

	// html输出
	Html(template string, obj interface{}) Response

	// string
	Text(format string, values ...interface{}) Response

	// 重定向
	Redirect(path string) Response

	// header
	SetHeader(key string, val string) Response

	// Cookie
	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) Response

	// 设置状态码
	SetStatus(code int) Response

	// 设置200状态
	SetOkStatus() Response
}

// Jsonp输出
func (ctx *Context) Jsonp(obj interface{}) Response {
	// 获取请求参数callback
	callbackFunc := ctx.Query("callback")
	ctx.SetHeader("Content-Type", "application/javascript")
	// 输出到前端页面的时候需要注意下进行字符过滤，否则有可能造成xss攻击
	callback := template.JSEscapeString(callbackFunc)

	// 输出函数名
	_, err := ctx.Writer.Write([]byte(callback))
	if err != nil {
		return ctx
	}
	// 输出左括号
	_, err = ctx.Writer.Write([]byte("("))
	if err != nil {
		return ctx
	}
	// 数据函数参数
	ret, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}
	_, err = ctx.Writer.Write(ret)
	if err != nil {
		return ctx
	}
	// 输出右括号
	_, err = ctx.Writer.Write([]byte(")"))
	if err != nil {
		return ctx
	}
	return ctx
}

// xml输出
func (ctx *Context) Xml(obj interface{}) Response {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/html")
	ctx.Writer.Write(byt)
	return ctx
}

// html输出
func (ctx *Context) Html(file string, obj interface{}) Response {
	// 读取模版文件，创建template实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx
	}
	// 执行Execute方法将obj和模版进行结合
	if err := t.Execute(ctx.Writer, obj); err != nil {
		return ctx
	}

	ctx.SetHeader("Content-Type", "application/html")
	return ctx
}

// string
func (ctx *Context) Text(format string, values ...interface{}) Response {
	out := fmt.Sprintf(format, values...)
	ctx.SetHeader("Content-Type", "application/json; charset=UTF-8")
	ctx.Writer.Write([]byte(out))
	return ctx
}

// 重定向
func (ctx *Context) Redirect(path string) Response {
	http.Redirect(ctx.Writer, ctx.Request, path, http.StatusMovedPermanently)
	return ctx
}

// header
func (ctx *Context) SetHeader(key string, val string) Response {
	ctx.Writer.Header().Add(key, val)
	return ctx
}

// Cookie
func (ctx *Context) SetCookie(
	key string,
	val string,
	maxAge int,
	path string,
	domain string,
	secure bool,
	httpOnly bool,
) Response {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return ctx
}

// 设置状态码
func (ctx *Context) SetStatus(code int) Response {
	ctx.Writer.WriteHeader(code)
	return ctx
}

// 设置200状态
func (ctx *Context) SetOkStatus() Response {
	ctx.Writer.WriteHeader(http.StatusOK)
	return ctx
}

func (ctx *Context) Json(obj interface{}) Response {
	byt, err := json.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Writer.Write(byt)
	return ctx
}
