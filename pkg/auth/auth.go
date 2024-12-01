package auth

import (
	"crypto/hmac"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	sha256 "github.com/minio/sha256-simd"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

const (
	TimestampDelay  = 30
	HeaderTimestamp = "x-timestamp"
	HeaderSignature = "x-signature"
	HeaderAuth      = "Authorization"
)

// RolePermissions maps roles to permission levels
var RolePermissions = map[string]int{
	"guest":      1,
	"user":       2,
	"premium":    5,
	"admin":      10,
	"device":     3,
	"superadmin": 15,
}

type CustomClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type Authenticator struct {
	deviceToken string
	jwtSecret   []byte
}

func NewAuthenticator(deviceToken string, jwtSecret string) *Authenticator {
	return &Authenticator{
		deviceToken: deviceToken,
		jwtSecret:   []byte(jwtSecret),
	}
}

func (a *Authenticator) GenerateSignature(ts string) string {
	h := hmac.New(sha256.New, []byte(a.deviceToken))
	h.Write([]byte(ts))
	return hex.EncodeToString(h.Sum(nil))
}

func (a *Authenticator) ValidateDeviceRequest(timestamp, signature string) error {
	if timestamp == "" || signature == "" {
		return errors.New("missing required headers")
	}
	if a.deviceToken == "" {
		return errors.New("device token is not configured")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "cannot parse timestamp")
	}

	currentTime := time.Now().UTC().Unix()
	if currentTime-ts > TimestampDelay || ts > currentTime+TimestampDelay {
		return errors.New("timestamp expired or not yet valid")
	}

	expectedSignature := a.GenerateSignature(timestamp)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return errors.New("invalid signature")
	}

	return nil
}

func (a *Authenticator) ValidateJWTToken(tokenString string) (*CustomClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

func (a *Authenticator) GenerateToken(userID int, role string, expiresIn time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresIn).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

// GetPermissionLevel returns the permission level for a given role
func GetPermissionLevel(role string) int {
	if level, exists := RolePermissions[strings.ToLower(role)]; exists {
		return level
	}
	return RolePermissions["guest"]
}
