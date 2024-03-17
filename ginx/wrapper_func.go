package ginx

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/apus-run/sea-kit/ginx/internal/errs"
	"github.com/apus-run/sea-kit/ginx/session"
)

func W(fn func(ctx *Context) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(&Context{Context: ctx})
		if errors.Is(err, errs.ErrNoResponse) {
			slog.Debug("不需要响应", slog.Any("err", err))
			return
		}
		if errors.Is(err, errs.ErrUnauthorized) {
			slog.Debug("未授权", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("执行业务逻辑失败", slog.Any("err", err))
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func B[Req any](fn func(ctx *Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			slog.Debug("绑定参数失败", slog.Any("err", err))
			return
		}
		res, err := fn(&Context{Context: ctx}, req)
		if errors.Is(err, errs.ErrNoResponse) {
			slog.Debug("不需要响应", slog.Any("err", err))
			return
		}
		if errors.Is(err, errs.ErrUnauthorized) {
			slog.Debug("未授权", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("执行业务逻辑失败", slog.Any("err", err))
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

// BS 的意思是，传入的业务逻辑方法可以接受 req 和 sess 两个参数
func BS[Req any](fn func(ctx *Context, req Req, sess session.Session) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gtx := &Context{Context: ctx}
		sess, err := session.Get(gtx)
		if err != nil {
			slog.Debug("获取 Session 失败", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var req Req
		// Bind 方法本身会返回 400 的错误
		if err := ctx.Bind(&req); err != nil {
			slog.Debug("绑定参数失败", slog.Any("err", err))
			return
		}
		res, err := fn(gtx, req, sess)
		if errors.Is(err, errs.ErrNoResponse) {
			slog.Debug("不需要响应", slog.Any("err", err))
			return
		}
		// 如果里面有权限校验，那么会返回 401 错误（目前来看，主要是登录态校验）
		if errors.Is(err, errs.ErrUnauthorized) {
			slog.Debug("未授权", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("执行业务逻辑失败", slog.Any("err", err))
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

// S 的意思是，传入的业务逻辑方法可以接受 Session 参数
func S(fn func(ctx *Context, sess session.Session) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gtx := &Context{Context: ctx}
		sess, err := session.Get(gtx)
		if err != nil {
			slog.Debug("获取 Session 失败", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := fn(gtx, sess)
		if errors.Is(err, errs.ErrNoResponse) {
			slog.Debug("不需要响应", slog.Any("err", err))
			return
		}
		// 如果里面有权限校验，那么会返回 401 错误（目前来看，主要是登录态校验）
		if errors.Is(err, errs.ErrUnauthorized) {
			slog.Debug("未授权", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("执行业务逻辑失败", slog.Any("err", err))
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}
