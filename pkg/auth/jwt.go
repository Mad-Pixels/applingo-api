package auth

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// Claims represents JWT claims structure
type Claims struct {
	Identifier int  `json:"identifier"`
	Role       Role `json:"role"`
	jwt.RegisteredClaims
}

// JWTAuth handles JWT-specific authentication
type JWTAuth struct {
	secret []byte
}

// NewJWTAuth creates new JWT authenticator instance
func NewJWTAuth(secret string) *JWTAuth {
	return &JWTAuth{
		secret: []byte(secret),
	}
}

// ValidateToken validates JWT token and returns claims
func (j *JWTAuth) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(strings.TrimPrefix(tokenString, "Bearer "), &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return j.secret, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidTokenClaims
}

// GenerateToken creates new JWT token with provided claims
func (j *JWTAuth) GenerateToken(identifier int, role Role, expiresIn time.Duration) (string, error) {
	claims := Claims{
		Identifier: identifier,
		Role:       role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}
