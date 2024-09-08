package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"os"
)

const (
	EnvLogLevel = "LOG_LEVEL"
)

// BaseRequest ...
type BaseRequest struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

// BaseResponse ...
type BaseResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
}

// HandleFunc ...
type HandleFunc func(context.Context, json.RawMessage) (any, error)

// Lambda ...
type Lambda struct {
	handlers map[string]HandleFunc
	logLvl   string
}

// MustLambda ...
func MustLambda(handlers map[string]HandleFunc) *Lambda {
	return &Lambda{
		logLvl:   os.Getenv(EnvLogLevel),
		handlers: handlers,
	}
}

// Route ...
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
		return errResponse(500, fmt.Sprintf("Error processing request: %v", err))
	}
	return okResponse(result)
}
