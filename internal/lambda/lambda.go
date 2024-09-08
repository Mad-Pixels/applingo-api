package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

const (
	EnvLogLevel = "LOG_LEVEL"
)

// BaseRequest represents the structure of incoming requests.
type BaseRequest struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

// HandleFunc is the type for action handlers.
type HandleFunc func(context.Context, json.RawMessage) (interface{}, error)

type lambda struct {
	handlers map[string]HandleFunc
	logger   *zap.Logger
}

// NewLambda creates a new Lambda object.
func NewLambda(handlers map[string]HandleFunc) (*lambda, error) {
	logger, err := initLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}
	return &lambda{
		handlers: handlers,
		logger:   logger,
	}, nil
}

// Handle processes the incoming Lambda event
func (l *lambda) Handle(ctx context.Context, event json.RawMessage) (events.APIGatewayProxyResponse, error) {
	l.logger.Info("Received event", zap.String("event", string(event)))

	var base BaseRequest
	if err := json.Unmarshal(event, &base); err != nil {
		l.logger.Error("Invalid request format", zap.Error(err))
		return NewResponse[interface{}](400, fmt.Sprintf("Invalid request format: %v", err), nil).ToAPIGatewayProxyResponse()
	}

	handler, ok := l.handlers[base.Action]
	if !ok {
		l.logger.Warn("Unknown action", zap.String("action", base.Action))
		return NewResponse[interface{}](404, fmt.Sprintf("Unknown action: %s", base.Action), nil).ToAPIGatewayProxyResponse()
	}

	result, err := handler(ctx, base.Data)
	if err != nil {
		l.logger.Error("Error processing request", zap.Error(err))
		return NewResponse[interface{}](500, fmt.Sprintf("Error processing request: %v", err), nil).ToAPIGatewayProxyResponse()
	}

	l.logger.Info("Request processed successfully", zap.Any("result", result))
	return NewResponse(200, "", result).ToAPIGatewayProxyResponse()
}
