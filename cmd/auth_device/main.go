package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

const (
	timestampDelay = 60

	headerTimestamp = "X-Timestamp"
	headerSignature = "X-Signature"
)

var (
	token  = os.Getenv("AUTH_TOKEN")
	logger = lambda.InitLogger()
)

func validateRequest(req events.APIGatewayCustomAuthorizerRequestTypeRequest, token string) error {
	var (
		timestamp = req.Headers[headerTimestamp]
		signature = req.Headers[headerSignature]

		generateSignature = func(ts, arn, t string) string {
			h := hmac.New(sha256.New, []byte(t))
			h.Write([]byte(ts))
			h.Write([]byte(arn))
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
		return errors.Wrap(err, "cannot parse timestamp")
	}
	currentTime := time.Now().UTC().Unix()
	if currentTime-ts > timestampDelay || ts > currentTime+timestampDelay {
		return errors.New("timestamp expired or not yet valid")
	}
	if !hmac.Equal([]byte(signature), []byte(generateSignature(timestamp, req.MethodArn, token))) {
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
		logger.Error().Err(err).Msg("Access denied")
		return generatePolicy("", "Deny", req.MethodArn)
	}
	return generatePolicy("device", "Allow", req.MethodArn)
}

func main() { aws_lambda.Start(handler) }
