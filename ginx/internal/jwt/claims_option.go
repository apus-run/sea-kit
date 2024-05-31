package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Options struct {
	Expire        time.Duration     // 有效期
	EncryptionKey string            // 加密密钥
	DecryptKey    string            // 解密密钥
	Method        jwt.SigningMethod // 签名方式
	Issuer        string            // 签发人
	genIDFn       func() string     // 生成 JWT ID (jti) 的函数
}

// NewOptions 定义一个 JWT 配置.
// DecryptKey: 默认与 EncryptionKey 相同.
// Method: 默认使用 jwt.SigningMethodHS256 签名方式.
func NewOptions(expire time.Duration, encryptionKey string,
	opts ...Option[Options]) Options {
	dOpts := Options{
		Expire:        expire,
		EncryptionKey: encryptionKey,
		DecryptKey:    encryptionKey,
		Method:        jwt.SigningMethodHS256,
		genIDFn:       func() string { return "" },
	}

	Apply[Options](&dOpts, opts...)

	return dOpts
}

// WithDecryptKey 设置解密密钥.
func WithDecryptKey(decryptKey string) Option[Options] {
	return func(o *Options) {
		o.DecryptKey = decryptKey
	}
}

// WithMethod 设置 JWT 的签名方法.
func WithMethod(method jwt.SigningMethod) Option[Options] {
	return func(o *Options) {
		o.Method = method
	}
}

// WithIssuer 设置签发人.
func WithIssuer(issuer string) Option[Options] {
	return func(o *Options) {
		o.Issuer = issuer
	}
}

// WithGenIDFunc 设置生成 JWT ID 的函数.
// 可以设置成 WithGenIDFunc(uuid.NewString).
func WithGenIDFunc(fn func() string) Option[Options] {
	return func(o *Options) {
		o.genIDFn = fn
	}
}
