package main

import (
	"context"
	"os"

	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/Mad-Pixels/applingo-api/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

var (
	deviceToken = os.Getenv("DEVICE_API_TOKEN")
	jwtSecret   = os.Getenv("JWT_SECRET")

	log           = logger.InitLogger()
	authenticator *auth.Authenticator
)

func init() {
	if deviceToken == "" || jwtSecret == "" {
		log.Fatal().Msg("AUTH_TOKEN and JWT_SECRET environment variables must be set")
	}
	authenticator = auth.NewAuthenticator(deviceToken, jwtSecret)
}

func generatePolicy(principalID string, effect string, resource string, context map[string]interface{}) (events.APIGatewayCustomAuthorizerResponse, error) {
	if effect != "Allow" && effect != "Deny" {
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("invalid effect")
	}
	authResponse := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: principalID,
		Context:     context,
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
	return authResponse, nil
}

func handler(ctx context.Context, req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	switch {
	case req.Headers[auth.HeaderTimestamp] != "" && req.Headers[auth.HeaderSignature] != "":
		return handleDeviceAuth(req)
	case req.Headers[auth.HeaderAuth] != "":
		return handleUserAuth(req)
	default:
		log.Error().Msg("authorization failed: No valid authentication headers found")
		return generatePolicy("", "Deny", req.MethodArn, nil)
	}
}

func main() {
	lambda.Start(handler)
}
