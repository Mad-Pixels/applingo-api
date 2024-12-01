package main

import (
	"context"
	"os"
	"strconv"

	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

var authenticator *auth.Authenticator

func init() {
	deviceToken := os.Getenv("AUTH_TOKEN")
	jwtSecret := os.Getenv("JWT_SECRET")
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
	// Device authentication
	if timestamp, ok := req.Headers[auth.HeaderTimestamp]; ok {
		if signature, ok := req.Headers[auth.HeaderSignature]; ok {
			if err := authenticator.ValidateDeviceRequest(timestamp, signature); err != nil {
				return generatePolicy("", "Deny", req.MethodArn, nil)
			}
			context := map[string]interface{}{
				"permissions": auth.RolePermissions["device"],
				"auth_type":   "device",
			}
			return generatePolicy("device", "Allow", req.MethodArn, context)
		}
	}

	// JWT authentication
	if authHeader, ok := req.Headers[auth.HeaderAuth]; ok {
		claims, err := authenticator.ValidateJWTToken(authHeader)
		if err != nil {
			return generatePolicy("", "Deny", req.MethodArn, nil)
		}

		permLevel := auth.GetPermissionLevel(claims.Role)
		context := map[string]interface{}{
			"user_id":     claims.UserID,
			"permissions": permLevel,
			"role":        claims.Role,
			"auth_type":   "jwt",
		}

		return generatePolicy(strconv.Itoa(claims.UserID), "Allow", req.MethodArn, context)
	}

	return generatePolicy("", "Deny", req.MethodArn, nil)
}

func main() {
	lambda.Start(handler)
}
