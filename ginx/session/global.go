package session

import (
	"github.com/gin-gonic/gin"

	"github.com/apus-run/sea-kit/ginx"
)

const CtxSessionKey = "_session"

var defaultProvider Provider

func NewSession(ctx *ginx.Context, uid int64,
	jwtData map[string]string,
	sessData map[string]any) (Session, error) {
	return defaultProvider.NewSession(
		ctx,
		uid,
		jwtData,
		sessData)
}

// Get 参考 defaultProvider.Get 的说明
func Get(ctx *ginx.Context) (Session, error) {
	return defaultProvider.Get(ctx)
}

func SetDefaultProvider(sp Provider) {
	defaultProvider = sp
}

func CheckLoginMiddleware() gin.HandlerFunc {
	return (&MiddlewareBuilder{sp: defaultProvider}).Build()
}

func RenewAccessToken(ctx *ginx.Context) error {
	return defaultProvider.RenewAccessToken(ctx)
}
