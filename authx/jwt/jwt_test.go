package jwt

import (
	"context"
	"github.com/apus-run/sea-kit/authx/jwt/store/redis"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID uint64

	// UserAgent 增强安全性，防止token被盗用
	UserAgent string

	jwt.RegisteredClaims
}

func TestGenerateToken(t *testing.T) {

}

func TestParseToken(t *testing.T) {

}

func TestNewJwtAuth(t *testing.T) {
	headers := make(map[string]any)
	headers["kid"] = "8b5228a5-b3d2-4165-aaac-58a052629846"
	now := time.Now()
	expiresAt := now.Add(2 * time.Hour)
	opts := []Option{

		WithTokenHeader(headers),
		WithExpired(2 * time.Hour),
		WithKeyfunc(func(token *jwt.Token) (any, error) {
			// Verify that the signing method is HMAC.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrTokenInvalid
			}
			return []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"), nil
		}),
		WithClaims(func() jwt.Claims {
			return &CustomClaims{
				UserID:    1,
				UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.83 Safari/537.36",
				RegisteredClaims: jwt.RegisteredClaims{
					// Issuer = iss,令牌颁发者。它表示该令牌是由谁创建的
					Issuer: "",
					// IssuedAt = iat,令牌颁发时的时间戳。它表示令牌是何时被创建的
					IssuedAt: jwt.NewNumericDate(now),
					// ExpiresAt = exp,令牌的过期时间戳。它表示令牌将在何时过期
					ExpiresAt: jwt.NewNumericDate(expiresAt),
					// NotBefore = nbf,令牌的生效时的时间戳。它表示令牌从什么时候开始生效
					NotBefore: jwt.NewNumericDate(now),
					// Subject = sub,令牌的主体。它表示该令牌是关于谁的
					Subject: "",
				},
			}
		}),
	}

	opts = append(opts, WithSigningMethod(jwt.SigningMethodHS256))

	store := redis.NewStore(nil, "authx")

	j, err := NewJwtAuth(store, opts...).Sign(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(j.GetToken())
}
