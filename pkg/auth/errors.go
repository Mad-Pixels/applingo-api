package auth

import "github.com/pkg/errors"

var (
	// Device authentication errors
	ErrMissingHeaders   = errors.New("missing required headers")
	ErrNoDeviceToken    = errors.New("device token is not configured")
	ErrTimestampParse   = errors.New("cannot parse timestamp")
	ErrTimestampExpired = errors.New("timestamp expired or not yet valid")
	ErrInvalidSignature = errors.New("invalid signature")

	// JWT authentication errors
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
)
