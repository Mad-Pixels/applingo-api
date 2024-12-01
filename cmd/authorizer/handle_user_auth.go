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

	permLevel := auth.GetPermissionLevel(claims.Role)
	context := map[string]interface{}{
		"user_id":     claims.Identifier,
		"permissions": permLevel,
		"role":        auth.RoleNames[claims.Role],
		"auth_type":   auth.JWT,
	}
	return generatePolicy(strconv.Itoa(claims.Identifier), "Allow", req.MethodArn, context)
}
