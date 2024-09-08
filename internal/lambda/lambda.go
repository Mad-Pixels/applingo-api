package lambda

import (
	"context"
	"encoding/json"
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
func NewLambda(handlers map[string]HandleFunc) *lambda {
	logger, _ := initLogger()

	return &lambda{
		handlers: handlers,
		logger:   logger,
	}
}

// Handle processes the incoming Lambda event
func (l *lambda) Handle(ctx context.Context, event json.RawMessage) (resp events.APIGatewayProxyResponse, err error) {
	l.logger.Debug("Received event", zap.String("event", string(event)))

	var (
		base   BaseRequest
		logMsg string
	)
	defer func() {
		if err != nil {
			l.logger.Error(logMsg, zap.Error(err), zap.Int("statusCode", resp.StatusCode), zap.String("action", base.Action))
		} else {
			l.logger.Info(logMsg, zap.Int("statusCode", resp.StatusCode), zap.String("action", base.Action))
		}
	}()

	if err = json.Unmarshal(event, &base); err != nil {
		logMsg = "Invalid request format"
		return NewResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil).ToAPIGatewayProxyResponse()
	}
	handler, ok := l.handlers[base.Action]
	if !ok {
		logMsg = "Unknown action"
		return NewResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil).ToAPIGatewayProxyResponse()
	}
	result, err := handler(ctx, base.Data)
	if err != nil {
		logMsg = "Error processing request"
		return NewResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil).ToAPIGatewayProxyResponse()
	}
	logMsg = "Request processed successfully"
	return NewResponse(http.StatusOK, http.StatusText(http.StatusOK), result).ToAPIGatewayProxyResponse()
}
