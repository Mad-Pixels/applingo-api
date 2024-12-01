package main

import (
	"strconv"

	"github.com/Mad-Pixels/applingo-api/pkg/auth"

	"github.com/aws/aws-lambda-go/events"
)

func handleUserAuth(req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	authHeader, ok := req.Headers[auth.HeaderAuth]
	if !ok {
		log.Error().Msg("authorization header missing for JWT authentication")
		return generatePolicy("", "Deny", req.MethodArn, nil)
	}
	claims, err := authenticator.ValidateJWTToken(authHeader)
	if err != nil {
		log.Error().Err(err).Msg("JWT authentication failed")
		return generatePolicy("", "Deny", req.MethodArn, nil)
	}
	context := map[string]interface{}{
		"identifier":  claims.Identifier,
		"permissions": strconv.Itoa(auth.GetPermissionLevel(auth.Device)),
		"role":        strconv.Itoa(int(auth.Device)),
		"kind":        strconv.Itoa(int(auth.HMAC)),
	}
	return generatePolicy(strconv.Itoa(claims.Identifier), "Allow", req.MethodArn, context)
}
