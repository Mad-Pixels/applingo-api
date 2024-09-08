package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
	"net/http"
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
type HandleFunc func(context.Context, *zap.Logger, json.RawMessage) (any, *HandleError)

// HandleError implement error from handlers.
type HandleError struct {
	Err    error
	Status int
}

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
		base BaseRequest
	)
	defer func() {
		if err != nil {
			l.logger.Error("Prepare response error", zap.Error(err), zap.String("action", base.Action))
			err = nil
		}
	}()
	if err = json.Unmarshal(event, &base); err != nil {
		l.logger.Error("Invalid request format", zap.Error(err), zap.String("action", base.Action))
		return errResponse(http.StatusInternalServerError)
	}

	handler, ok := l.handlers[base.Action]
	if !ok {
		l.logger.Error("Unknown action", zap.Error(errors.New("requested action not implemented")), zap.String("action", base.Action))
		return errResponse(http.StatusNotFound)
	}
	result, handleError := handler(ctx, l.logger, base.Data)
	if handleError != nil {
		l.logger.Error("Handle error", zap.Error(handleError.Err), zap.String("action", base.Action))
		return errResponse(handleError.Status)
	}
	return okResponse(result)
}
