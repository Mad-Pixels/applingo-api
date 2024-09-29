package api

import (
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/aws/aws-lambda-go/events"
)

func response(statusCode int, body any) (events.APIGatewayProxyResponse, error) {
	jsonBody, err := serializer.MarshalJSON(body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonBody),
		StatusCode: statusCode,
	}, nil
}

func errResponse(statusCode int) (events.APIGatewayProxyResponse, error) {
	return response(statusCode, map[string]string{"error": http.StatusText(statusCode)})
}

func okResponse(data any) (events.APIGatewayProxyResponse, error) {
	return response(http.StatusOK, map[string]any{"data": data})
}
