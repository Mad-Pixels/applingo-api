package main

import (
	"strconv"

	"github.com/Mad-Pixels/applingo-api/pkg/auth"

	"github.com/aws/aws-lambda-go/events"
)

func handleDeviceAuth(req events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	timestamp, tsOK := req.Headers[auth.HeaderTimestamp]
	signature, sigOK := req.Headers[auth.HeaderSignature]

	if !tsOK || !sigOK {
		log.Error().Msg("missing required headers for device authentication")
		return generatePolicy("", "Deny", req.MethodArn, nil)
	}
	if err := authenticator.ValidateDeviceRequest(timestamp, signature); err != nil {
		log.Error().Err(err).Msg("device authentication failed")
		return generatePolicy("", "Deny", req.MethodArn, nil)
	}
	context := map[string]interface{}{
		"permissions": strconv.Itoa(auth.GetPermissionLevel(auth.Device)),
		"role":        strconv.Itoa(int(auth.Device)),
		"kind":        strconv.Itoa(int(auth.HMAC)),
	}
	return generatePolicy("device", "Allow", req.MethodArn, context)
}
