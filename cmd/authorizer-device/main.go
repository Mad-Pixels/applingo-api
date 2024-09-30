package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/Mad-Pixels/lingocards-api/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

const (
	timestampDelay  = 30
	headerTimestamp = "x-timestamp"
	headerSignature = "x-signature"
)

var (
	token = os.Getenv("AUTH_TOKEN")
	log   = logger.InitLogger()
)

func init() {
	debug.SetGCPercent(500)
}

func generateSignature(ts, token string) string {
	h := hmac.New(sha256.New, []byte(token))
	h.Write([]byte(ts))
	return hex.EncodeToString(h.Sum(nil))
}

func validateRequest(req events.APIGatewayCustomAuthorizerRequestTypeRequest, token string) error {
	var (
		timestamp = req.Headers[headerTimestamp]
		signature = req.Headers[headerSignature]
	)
	if timestamp == "" || signature == "" {
		return errors.New("missing required headers")
	}
	if token == "" {
		return errors.New("token is empty")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "cannot parse timestamp")
	}
	currentTime := time.Now().UTC().Unix()
	if currentTime-ts > timestampDelay || ts > currentTime+timestampDelay {
		return errors.New("timestamp expired or not yet valid")
	}
	expectedSignature := generateSignature(timestamp, token)

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return errors.New("invalid signature")
	}
	return nil
}

func generatePolicy(principalID, effect, resource string) (events.APIGatewayCustomAuthorizerResponse, error) {
	if effect != "Allow" && effect != "Deny" {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("invalid effect")
	}
	policy := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: principalID,
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		},
	}
	return policy, nil
}

func handler(_ context.Context, req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	if err := validateRequest(req, token); err != nil {
		log.Error().Err(err).Msg("Access denied")
		return generatePolicy("", "Deny", req.MethodArn)
	}
	return generatePolicy("device", "Allow", req.MethodArn)
}

func main() { lambda.Start(handler) }
