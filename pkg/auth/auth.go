// Package auth defines user roles and permissions used for access control.
package auth

import "time"

const (
	// TimestampDelay defines the allowed time difference (in seconds) for validating request timestamps.
	TimestampDelay = 15

	// HeaderTimestamp is the HTTP header name used to send a request timestamp.
	HeaderTimestamp = "x-timestamp"

	// HeaderSignature is the HTTP header name used to send the HMAC signature.
	HeaderSignature = "x-signature"

	// HeaderAuth is the HTTP header name used to send the JWT token.
	HeaderAuth = "Authorization"
)

// Authenticator provides the main authentication functionality
// combining HMAC and JWT-based mechanisms.
type Authenticator struct {
	deviceToken string
	jwtSecret   []byte
	hmac        *HMACAuth
	jwt         *JWTAuth
}

// NewAuthenticator creates a new instance of Authenticator.
func NewAuthenticator(deviceToken string, jwtSecret string) *Authenticator {
	auth := &Authenticator{
		deviceToken: deviceToken,
		jwtSecret:   []byte(jwtSecret),
	}
	auth.hmac = NewHMACAuth(deviceToken)
	auth.jwt = NewJWTAuth(jwtSecret)
	return auth
}

// ValidateDeviceRequest validates device authentication request using timestamp and signature.
func (a *Authenticator) ValidateDeviceRequest(timestamp, signature string) error {
	return a.hmac.ValidateRequest(timestamp, signature)
}

// ValidateJWTToken validates a JWT token and returns the parsed claims.
func (a *Authenticator) ValidateJWTToken(tokenString string) (*Claims, error) {
	return a.jwt.ValidateToken(tokenString)
}

// GenerateToken generates a new JWT token with the given user ID, role, and expiration time.
func (a *Authenticator) GenerateToken(userID int, role Role, expiresIn time.Duration) (string, error) {
	return a.jwt.GenerateToken(userID, role, expiresIn)
}
