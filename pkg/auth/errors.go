package auth

import "github.com/pkg/errors"

var (
	// ErrMissingHeaders indicates that required headers are missing from the request.
	ErrMissingHeaders = errors.New("missing required headers")

	// ErrNoDeviceToken indicates that the device token is not configured.
	ErrNoDeviceToken = errors.New("device token is not configured")

	// ErrTimestampParse indicates a failure to parse the provided timestamp.
	ErrTimestampParse = errors.New("cannot parse timestamp")

	// ErrTimestampExpired means the provided timestamp is either expired or not yet valid.
	ErrTimestampExpired = errors.New("timestamp expired or not yet valid")

	// ErrInvalidSignature means the HMAC signature validation failed.
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrUnexpectedSigningMethod means the JWT signing method does not match the expected one.
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")

	// ErrInvalidTokenClaims means the JWT token contains invalid claims.
	ErrInvalidTokenClaims = errors.New("invalid token claims")
)
