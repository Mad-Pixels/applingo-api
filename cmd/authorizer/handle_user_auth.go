package main

import (
	"strconv"

	"github.com/Mad-Pixels/applingo-api/pkg/auth"

	"github.com/aws/aws-lambda-go/events"
)

func handleUserAuth(token string, req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	claims, err := authenticator.ValidateJWTToken(token)
	if err != nil {
		log.Error().Err(err).Msg("JWT authentication failed")
		return generatePolicy("", "Deny", req.MethodArn, nil)
	}
	context := map[string]interface{}{
		"identifier":  claims.Identifier,
		"permissions": strconv.Itoa(auth.GetPermissionLevel(auth.User)),
		"role":        strconv.Itoa(int(auth.User)),
		"kind":        strconv.Itoa(int(auth.JWT)),
	}
	return generatePolicy(strconv.Itoa(claims.Identifier), "Allow", req.MethodArn, context)
}
