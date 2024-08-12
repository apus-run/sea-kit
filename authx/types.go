package authx

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Token interface {
	// GetToken Get token string.
	GetToken() string
	// GetTokenType Get token type.
	GetTokenType() string
	// GetExpireAt Get token expiration timestamp.
	GetExpireAt() int64
	// EncodeToJSON JSON encoding
	EncodeToJSON() ([]byte, error)
}

// Authenticator defines methods used for token processing.
type Authenticator interface {
	// Sign is used to generate a token.
	Sign(ctx context.Context, userID string) (Token, error)

	// Destroy is used to destroy a token.
	Destroy(ctx context.Context, accessToken string) error

	// ParseClaims parse the token and return the claims.
	ParseClaims(ctx context.Context, accessToken string) (*jwt.RegisteredClaims, error)

	// ParseToken is used to parse a token.
	ParseToken(ctx context.Context, accessToken string) (*jwt.Token, error)

	// GenerateToken is used to generate a token.
	GenerateToken(ctx context.Context) (string, error)
}

// Encrypt encrypts the plain text with bcrypt.
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare compares the encrypted text with the plain text if it's the same.
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
