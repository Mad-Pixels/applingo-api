package lambda

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

func response(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonBody),
	}, nil
}

func errResponse(statusCode int) (events.APIGatewayProxyResponse, error) {
	body := map[string]string{
		"error": http.StatusText(statusCode),
	}
	return response(statusCode, body)
}

func okResponse(data interface{}) (events.APIGatewayProxyResponse, error) {
	body := map[string]interface{}{
		"data": data,
	}
	return response(http.StatusOK, body)
}

////go:generate msgp
