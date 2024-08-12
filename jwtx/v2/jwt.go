package jwtx

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingKeyFunc         = errors.New("keyFunc is missing")
	ErrTokenInvalid           = errors.New("token is invalid")
	ErrUnSupportSigningMethod = errors.New("wrong signing method")
	ErrNeedTokenProvider      = errors.New("token provider is missing")
	ErrSignToken              = errors.New("can not sign token. is the key correct")
	ErrGetKey                 = errors.New("can not get key while signing token")
)

// GenerateToken 生成jwt token
func GenerateToken(keyProvider jwt.Keyfunc, opts ...Option) (string, error) {
	o := Apply(opts...)

	if keyProvider == nil {
		return "", ErrNeedTokenProvider
	}
	token := jwt.NewWithClaims(o.signingMethod, o.claims())
	if o.tokenHeader != nil {
		for k, v := range o.tokenHeader {
			token.Header[k] = v
		}
	}
	key, err := keyProvider(token)
	if err != nil {
		return "", ErrGetKey
	}
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", ErrSignToken
	}

	return tokenStr, nil
}

// ParseToken 解析jwt token
func ParseToken(jwtToken string, keyFunc jwt.Keyfunc, opts ...Option) (token *jwt.Token, err error) {
	o := Apply(opts...)
	if keyFunc == nil {
		return nil, ErrMissingKeyFunc
	}

	if o.claims != nil {
		token, err = jwt.ParseWithClaims(jwtToken, o.claims(), keyFunc)
	} else {
		token, err = jwt.Parse(jwtToken, keyFunc)
	}

	// 过期的, 伪造的, 都可以认为是无效token
	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid
	}

	if token.Method != o.signingMethod {
		return nil, ErrUnSupportSigningMethod
	}

	return token, nil
}
