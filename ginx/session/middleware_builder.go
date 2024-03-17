package session

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/apus-run/sea-kit/ginx"
)

// MiddlewareBuilder 登录校验
type MiddlewareBuilder struct {
	sp Provider
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sess, err := b.sp.Get(&ginx.Context{Context: ctx})
		if err != nil {
			slog.Debug("未授权", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set(CtxSessionKey, sess)
	}
}
