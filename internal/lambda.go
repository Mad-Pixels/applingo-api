package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

type BaseRequest struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type BaseResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
}

type HandleFunc func(context.Context, json.RawMessage) (any, error)

type Lambda struct {
	AwsSes *session.Session

	handlers map[string]HandleFunc
	logLvl   string
}

// MustLambda ...
func MustLambda(handlers map[string]HandleFunc) *Lambda {
	return &Lambda{
		AwsSes: session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
		})),
		logLvl:   os.Getenv("LOG_LEVEL"),
		handlers: handlers,
	}
}

func (l *Lambda) Route(ctx context.Context, req json.RawMessage) (events.APIGatewayProxyResponse, error) {
	var base BaseRequest
	if err := json.Unmarshal(req, &base); err != nil {
		return errResponse(400, fmt.Sprintf("Invalid request format: %v", err))
	}

	handler, ok := l.handlers[base.Action]
	if !ok {
		return errResponse(404, fmt.Sprintf("Unknown action: %s", base.Action))
	}
	result, err := handler(ctx, base.Data)
	if err != nil {
		return errResponse(500, fmt.Sprintf("Error procissing request: %v", err))
	}
	return okResponse(result)
}

func errResponse(status int, message string) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(BaseResponse{
		StatusCode: status,
		Message:    message,
	})
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func okResponse(data any) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(BaseResponse{
		StatusCode: 200,
		Data:       data,
	})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}, nil
}
