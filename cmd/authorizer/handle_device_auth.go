package main

import (
	"strconv"

	"github.com/Mad-Pixels/applingo-api/pkg/auth"

	"github.com/aws/aws-lambda-go/events"
)

func handleDeviceAuth(timestamp string, signature string, req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	if err := authenticator.ValidateDeviceRequest(timestamp, signature); err != nil {
		log.Error().Err(err).Msg("Device authentication failed")
		return generatePolicy("", "Deny", req.MethodArn, nil)
	}
	context := map[string]interface{}{
		"permissions": strconv.Itoa(auth.GetPermissionLevel(auth.Device)),
		"role":        strconv.Itoa(int(auth.Device)),
		"kind":        strconv.Itoa(int(auth.HMAC)),
	}
	return generatePolicy("device", "Allow", req.MethodArn, context)
}
