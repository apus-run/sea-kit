package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/apus-run/sea-kit/authx"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenInvalid           = errors.New("token is invalid")
	ErrUnSupportSigningMethod = errors.New("wrong signing method")
	ErrSignToken              = errors.New("can not sign token. is the key correct")
	ErrGetKey                 = errors.New("can not get key while signing token")
)

// JwtAuth implement the authx.Authenticator interface.

type JwtAuth struct {
	*options
	store Storer
}

func NewJwtAuth(store Storer, opts ...Option) *JwtAuth {
	options := Apply(opts...)
	return &JwtAuth{
		options: options,
		store:   store,
	}
}

func (j *JwtAuth) Sign(ctx context.Context) (authx.Token, error) {
	now := time.Now()
	expiresAt := now.Add(j.expired)

	tokenString, err := j.GenerateToken(ctx)
	if err != nil {
		return nil, err
	}
	tokenInfo := &tokenInfo{
		Token:     tokenString,
		Type:      j.tokenType,
		ExpiresAt: expiresAt.Unix(),
	}
	return tokenInfo, nil
}

func (j *JwtAuth) Destroy(ctx context.Context, refreshToken string) error {
	claims, err := j.ParseClaims(ctx, refreshToken)
	if err != nil {
		return err
	}

	// If storage is set, put the unexpired token in
	store := func(store Storer) error {
		expired := time.Until(claims.ExpiresAt.Time)
		return store.Set(ctx, refreshToken, "1", expired)
	}
	return j.callStore(store)
}

func (j *JwtAuth) ParseClaims(ctx context.Context, accessToken string) (*jwt.RegisteredClaims, error) {
	if accessToken == "" {
		return nil, ErrTokenInvalid
	}

	token, err := j.ParseToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	store := func(store Storer) error {
		exists, err := store.Check(ctx, accessToken)
		if err != nil {
			return err
		}

		if exists {
			return ErrTokenInvalid
		}

		return nil
	}

	if err := j.callStore(store); err != nil {
		return nil, err
	}

	return token.Claims.(*jwt.RegisteredClaims), nil
}

func (j *JwtAuth) ParseToken(ctx context.Context, accessToken string) (token *jwt.Token, err error) {
	if j.claims != nil {
		token, err = jwt.ParseWithClaims(accessToken, j.claims(), j.keyfunc)
	} else {
		token, err = jwt.Parse(accessToken, j.keyfunc)
	}

	// 过期的, 伪造的, 都可以认为是无效token
	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid
	}

	if token.Method != j.signingMethod {
		return nil, ErrUnSupportSigningMethod
	}

	return token, nil
}

func (j *JwtAuth) GenerateToken(ctx context.Context) (string, error) {
	token := jwt.NewWithClaims(j.signingMethod, j.claims())
	if j.tokenHeader != nil {
		for k, v := range j.tokenHeader {
			token.Header[k] = v
		}
	}
	key, err := j.keyfunc(token)
	if err != nil {
		return "", ErrGetKey
	}
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", ErrSignToken
	}

	return tokenStr, nil
}

func (j *JwtAuth) callStore(fn func(Storer) error) error {
	if store := j.store; store != nil {
		return fn(store)
	}
	return nil
}
