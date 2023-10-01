package jwtx

import (
	"github.com/apus-run/sea-kit/log"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type CustomClaims struct {
	UserID uint64

	// UserAgent 增强安全性，防止token被盗用
	UserAgent string

	jwt.RegisteredClaims
}

func TestGenerateToken(t *testing.T) {
	testKey := "testKey"
	tProvider := func(*jwt.Token) (interface{}, error) {
		return []byte(testKey), nil
	}
	testCases := []struct {
		// 名字
		name string

		// 要提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after func(t *testing.T)

		// 预期入参
		token         func() string
		tokenProvider jwt.Keyfunc
		signingMethod jwt.SigningMethod
		claims        func() jwt.Claims
		key           string

		// 预期响应
		want    any
		wantErr error
	}{
		{
			name:   "成功生成token",
			before: func(t *testing.T) {},
			after:  func(t *testing.T) {},
			token: func() string {
				tokenStr, err := jwt.
					NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{}).
					SignedString([]byte(testKey))
				assert.NoError(t, err)
				return tokenStr
			},
			tokenProvider: tProvider,
			signingMethod: jwt.SigningMethodHS256,
			claims: func() jwt.Claims {
				return &CustomClaims{}
			},
			wantErr: nil,
			key:     testKey,
		},
		{
			name:   "CustomClaims 为 nil",
			before: func(t *testing.T) {},
			after:  func(t *testing.T) {},
			token: func() string {
				tokenStr, err := jwt.
					NewWithClaims(jwt.SigningMethodHS256, nil).
					SignedString([]byte(testKey))
				assert.NoError(t, err)
				return tokenStr
			},
			tokenProvider: tProvider,
			signingMethod: jwt.SigningMethodHS256,
			claims: func() jwt.Claims {
				return nil
			},
			wantErr: nil,
			key:     testKey,
		},
		{
			name:   "miss token provider",
			before: func(t *testing.T) {},
			after:  func(t *testing.T) {},
			token: func() string {
				tokenStr, err := jwt.
					NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{}).
					SignedString(nil)
				assert.Error(t, err)
				return tokenStr
			},
			tokenProvider: nil,
			signingMethod: jwt.SigningMethodHS512,
			claims: func() jwt.Claims {
				return &CustomClaims{}
			},
			wantErr: ErrNeedTokenProvider,
			key:     testKey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			tokenStr, err := GenerateToken(
				tc.tokenProvider,
				WithClaims(tc.claims),
				WithSigningMethod(jwt.SigningMethodHS256),
			)
			t.Logf("\n1: %s\n2: %s\n", tokenStr, tc.token())
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tokenStr, tc.token())
			tc.after(t)
		})
	}
}

func TestParseToken(t *testing.T) {
	testKey := "testKey"
	tProvider := func(*jwt.Token) (interface{}, error) {
		return []byte(testKey), nil
	}
	claims := &CustomClaims{}
	claims.UserID = 123
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(testKey))
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		// 名字
		name string

		// 要提前准备数据
		before func(t *testing.T)
		// 验证并且删除数据
		after func(t *testing.T)

		// 预期入参
		tokenStr      string
		tokenProvider jwt.Keyfunc
		signingMethod jwt.SigningMethod
		claims        func() jwt.Claims
		key           string

		// 预期响应
		want    any
		wantErr error
	}{
		{
			name:     "解析token成功",
			before:   func(t *testing.T) {},
			after:    func(t *testing.T) {},
			tokenStr: tokenStr,
			claims: func() jwt.Claims {
				return &CustomClaims{}
			},
			signingMethod: jwt.SigningMethodHS256,
			tokenProvider: tProvider,
			wantErr:       nil,
		},
		{
			name:     "miss key",
			before:   func(t *testing.T) {},
			after:    func(t *testing.T) {},
			tokenStr: tokenStr,
			claims: func() jwt.Claims {
				return &CustomClaims{}
			},
			signingMethod: jwt.SigningMethodHS256,
			tokenProvider: nil,
			wantErr:       ErrMissingKeyFunc,
		},
		{
			name:     "method invalid",
			before:   func(t *testing.T) {},
			after:    func(t *testing.T) {},
			tokenStr: tokenStr,
			claims: func() jwt.Claims {
				return &CustomClaims{}
			},
			signingMethod: jwt.SigningMethodHS512,
			tokenProvider: tProvider,
			wantErr:       ErrUnSupportSigningMethod,
		},
		{
			name:     "token invalid",
			before:   func(t *testing.T) {},
			after:    func(t *testing.T) {},
			tokenStr: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjAsIlVzZXJBZ2VudCI6IiJ9.xM8P27ckU6C3TxlW2fwPlFrr4P2ROE2hBoT3Gbls",
			claims: func() jwt.Claims {
				return &CustomClaims{}
			},
			signingMethod: jwt.SigningMethodHS512,
			tokenProvider: tProvider,
			wantErr:       ErrTokenInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			token, err := ParseToken(
				tc.tokenStr,
				tc.tokenProvider,
				WithClaims(tc.claims),
				WithSigningMethod(tc.signingMethod),
			)

			if err == nil {
				if c, ok := token.Claims.(*CustomClaims); ok && token.Valid {
					log.Infof("Claims: %v", c)
				}
			}

			assert.Equal(t, tc.wantErr, err)

			tc.after(t)
		})
	}
}
