package lambda

import (
	"encoding/json"
	"net/http"

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
func (r *response) ToAPIGatewayProxyResponse() events.APIGatewayProxyResponse {
	body, _ := json.Marshal(r)
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
	}
}

// NewResponse creates a new response object.
func NewResponse(statusCode int32, data any) *response {
	return &response{
		Headers:    make(map[string]string),
		StatusCode: statusCode,
		Message:    http.StatusText(int(statusCode)),
		Data:       data,
	}
}
