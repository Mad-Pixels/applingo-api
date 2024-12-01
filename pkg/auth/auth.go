package auth

import "time"

const (
	TimestampDelay  = 15
	HeaderTimestamp = "x-timestamp"
	HeaderSignature = "x-signature"
	HeaderAuth      = "Authorization"
)

// Authenticator provides the main authentication functionality
type Authenticator struct {
	deviceToken string
	jwtSecret   []byte
	hmac        *HMACAuth
	jwt         *JWTAuth
}

// NewAuthenticator creates a new instance of Authenticator
func NewAuthenticator(deviceToken string, jwtSecret string) *Authenticator {
	auth := &Authenticator{
		deviceToken: deviceToken,
		jwtSecret:   []byte(jwtSecret),
	}
	auth.hmac = NewHMACAuth(deviceToken)
	auth.jwt = NewJWTAuth(jwtSecret)
	return auth
}

// ValidateDeviceRequest validates device authentication request
func (a *Authenticator) ValidateDeviceRequest(timestamp, signature string) error {
	return a.hmac.ValidateRequest(timestamp, signature)
}

// ValidateJWTToken validates JWT token and returns claims
func (a *Authenticator) ValidateJWTToken(tokenString string) (*Claims, error) {
	return a.jwt.ValidateToken(tokenString)
}

// GenerateToken generates new JWT token
func (a *Authenticator) GenerateToken(userID int, role Role, expiresIn time.Duration) (string, error) {
	return a.jwt.GenerateToken(userID, role, expiresIn)
}
