package jwtx

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// SecretKey jwtx secret key
var SecretKey = "moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"

// CustomClaims 在标准声明中加入用户id
type CustomClaims struct {
	UserID uint64

	// UserAgent 增强安全性，防止token被盗用
	UserAgent string

	jwt.RegisteredClaims
}

// GenerateToken 生成jwt token
func GenerateToken(options ...Option) (string, error) {
	opts := Apply(options...)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID:    opts.userID,
		UserAgent: opts.userAgent,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(opts.expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})
	tokenString, err := token.SignedString([]byte(opts.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析jwt token
func ParseToken(tokenString, secretKey string) (*CustomClaims, *jwt.Token, error) {
	cc := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, cc, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return cc, nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, token, err
	}
	return cc, nil, err
}
