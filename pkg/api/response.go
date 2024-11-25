package api

import (
	"net/http"

	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/aws/aws-lambda-go/events"
)

func gatewayResponse(statusCode int, body any, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}

	jsonBody, err := serializer.MarshalJSON(body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       string(jsonBody),
	}, nil
}
