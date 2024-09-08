package lambda

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

//go:generate msgp

type response struct {
	Headers    map[string]string `json:"-" msg:"-"`
	Data       any               `json:"data,omitempty" msg:"data"`
	StatusCode int32             `json:"status_code" msg:"status_code"`
	Message    string            `json:"message,omitempty" msg:"message"`
}

// SetHeader sets a header for the response.
func (r *response) SetHeader(key, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = value
}

// ToAPIGatewayProxyResponse creates a new events.APIGatewayProxyResponse object.
func (r *response) ToAPIGatewayProxyResponse() (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	if _, ok := r.Headers["Content-Type"]; !ok {
		r.Headers["Content-Type"] = "application/json"
	}

	return events.APIGatewayProxyResponse{
		StatusCode: int(r.StatusCode),
		Headers:    r.Headers,
		Body:       string(body),
	}, nil
}

// NewResponse creates a new response object.
func NewResponse(statusCode int32, message string, data any) *response {
	return &response{
		Headers:    make(map[string]string),
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}
