package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// defaultKey holds the default key used to sign a jwt token.
	defaultKey = "authx::jwt(#)9527"
)

// Option is jwt option.
type Option func(*options)

// Parser is a jwt parser
type options struct {
	signingMethod jwt.SigningMethod
	claims        func() jwt.Claims
	tokenHeader   map[string]any

	expired   time.Duration
	keyfunc   jwt.Keyfunc
	tokenType string
}

// DefaultOptions .
func DefaultOptions() *options {
	return &options{
		tokenType:     "Bearer",
		expired:       2 * time.Hour,
		signingMethod: jwt.SigningMethodHS256,
		keyfunc: func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrTokenInvalid
			}
			return []byte(defaultKey), nil
		},
	}
}

func Apply(opts ...Option) *options {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// WithSigningMethod with signing method option.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// WithClaims with customer claim
// If you use it in Server, f needs to return a new jwt.Claims object each time to avoid concurrent write problems
// If you use it in Client, f only needs to return a single object to provide performance
func WithClaims(f func() jwt.Claims) Option {
	return func(o *options) {
		o.claims = f
	}
}

// WithTokenHeader withe customer tokenHeader for client side
func WithTokenHeader(header map[string]any) Option {
	return func(o *options) {
		o.tokenHeader = header
	}
}

// WithKeyfunc set the callback function for verifying the key.
func WithKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfunc = keyFunc
	}
}

// WithExpired set the token expiration time (in seconds, default 2h).
func WithExpired(expired time.Duration) Option {
	return func(o *options) {
		o.expired = expired
	}
}
