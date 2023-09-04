package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/apus-run/sea-kit/ginx"
	"github.com/apus-run/sea-kit/jwtx"
)

// Builder 鉴权，验证用户token是否有效
type Builder struct {
	// 白名单路由地址集合, 放行
	whitePathList []string
}

func NewBuilder() *Builder {
	return &Builder{
		whitePathList: []string{},
	}
}

func (b *Builder) IgnorePaths(whitePath string) *Builder {
	b.whitePathList = append(b.whitePathList, whitePath)
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return ginx.Handle(func(ctx *ginx.Context) {
		// 白名单路由放行
		for _, path := range b.whitePathList {
			if strings.Contains(ctx.Request.URL.Path, path) {
				ctx.Next()
				return
			}
		}

		tokenString, err := getJwtFromHeader(ctx)
		if err != nil {
			ctx.JSONE(http.StatusUnauthorized, "invalid token", nil)
			ctx.Abort()
			return
		}
		// 验证token是否正确
		claims, token, err := jwtx.ParseToken(tokenString, jwtx.SecretKey)
		if err != nil {
			ctx.JSONE(http.StatusUnauthorized, err.Error(), nil)
			ctx.Abort()
			return
		}

		// 验证token是否过期
		expireTime, err := claims.GetExpirationTime()
		if err != nil {
			// 拿不到过期时间
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		if expireTime.Before(now) {
			// 已经过期
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 利用 User-Agent 验证 token 是否被盗用
		if ctx.GetHeader("User-Agent") != claims.UserAgent {
			// User-Agent 不相等，可能是黑客伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 每 10 秒刷新一次
		if expireTime.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			newToken, err := token.SignedString(jwtx.SecretKey)
			if err != nil {
				// 因为刷新这个事情，并不是一定要做的，所以这里可以考虑打印日志
				// 暂时这样打印
				log.Println(err)

				return
			} else {
				ctx.Header("x-jwt-token", newToken)
				log.Println("refresh token")
			}
		}

		// 将claims信息存入上下文
		ctx.Set("claims", claims)
		ctx.Next()
	})
}

func getJwtFromHeader(ctx *ginx.Context) (string, error) {
	// 读取请求头的 token
	tokenString := ctx.GetHeader("Authorization")
	if len(tokenString) == 0 {
		return "", errors.New("token 为空")
	}
	strs := strings.SplitN(tokenString, " ", 2)
	if len(strs) != 2 || strs[0] != "Bearer" {
		return "", errors.New("token 不符合规则, Bearer 开头")
	}
	return strs[1], nil
}
