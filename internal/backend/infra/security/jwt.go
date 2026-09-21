// Package token contains the JWT adapter for the TokenIssuer port.
package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LeBaoTai/SDN-RD/internal/backend/domain/user"
	"github.com/golang-jwt/jwt/v5"
)

// Claims extends jwt.RegisteredClaims with our custom fields.
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// JWTIssuer implements user.TokenIssuer using HS256 signed JWTs.
type JWTIssuer struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTIssuer(secret string, ttl time.Duration) *JWTIssuer {
	return &JWTIssuer{secret: []byte(secret), ttl: ttl}
}

func (j *JWTIssuer) Issue(_ context.Context, u *user.User) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		Email: u.Email().String(),
		Role:  u.Role().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.ID().String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTIssuer) Parse(tokenStr string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return j.secret, nil
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, jwt.ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, jwt.ErrTokenNotValidYet
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return nil, jwt.ErrTokenSignatureInvalid
		default:
			return nil, err
		}
	}

	claims, ok := t.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
