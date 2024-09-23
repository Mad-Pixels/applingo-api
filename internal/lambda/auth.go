package lambda

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

const timestampDelay = 60

func validateRequest(req events.APIGatewayProxyRequest, token string) error {
	var (
		timestamp = req.Headers["X-Timestamp"]
		signature = req.Headers["X-Signature"]

		generateSignature = func(ts, p, t string) string {
			h := hmac.New(sha256.New, []byte(t))
			h.Write([]byte(ts))
			h.Write([]byte(p))
			return hex.EncodeToString(h.Sum(nil))
		}
	)
	if timestamp == "" || signature == "" {
		return errors.New("missing required headers")
	}
	if token == "" {
		return errors.New("token is empty")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return err
	}
	currentTime := time.Now().UTC().Unix()
	if currentTime-ts > timestampDelay || ts > currentTime+timestampDelay {
		return errors.New("timestamp expired or not yet valid")
	}

	if !hmac.Equal([]byte(signature), []byte(generateSignature(timestamp, req.Path, token))) {
		return errors.New("invalid signature")
	}
	return nil
}
