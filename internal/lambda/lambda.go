package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

const (
	EnvLogLevel = "LOG_LEVEL"
)

// BaseRequest represents the structure of incoming requests.
type BaseRequest struct {
	Data   json.RawMessage `json:"data"`
	Action string          `json:"action"`
}

// HandleFunc is the type for action handlers.
type HandleFunc func(context.Context, json.RawMessage) (any, error)

type lambda struct {
	logger   *zap.Logger
	handlers map[string]HandleFunc
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
	l.logger.Debug("Received event", zap.String("event", string(event)))

	var base BaseRequest
	if err := json.Unmarshal(event, &base); err != nil {
		l.logger.Error("Invalid request format", zap.Error(err))
		return NewResponse(400, http.StatusText(400), nil).ToAPIGatewayProxyResponse()
	}

	handler, ok := l.handlers[base.Action]
	if !ok {
		l.logger.Warn("Unknown action", zap.String("action", base.Action))
		return NewResponse(404, http.StatusText(404), nil).ToAPIGatewayProxyResponse()
	}

	result, err := handler(ctx, base.Data)
	if err != nil {
		l.logger.Error("Error processing request", zap.Error(err))
		return NewResponse(500, http.StatusText(500), nil).ToAPIGatewayProxyResponse()
	}

	l.logger.Debug("Request processed successfully", zap.Any("result", result))
	return NewResponse(200, http.StatusText(200), result).ToAPIGatewayProxyResponse()
}
