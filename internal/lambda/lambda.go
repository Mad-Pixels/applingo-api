package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"net/http"
)

// BaseRequest represents the structure of incoming requests.
type BaseRequest struct {
	Data   json.RawMessage `json:"data"`
	Action string          `json:"action"`
}

// HandleFunc is the type for action handlers.
type HandleFunc func(context.Context, zerolog.Logger, json.RawMessage) (any, *HandleError)

// HandleError implement error from handlers.
type HandleError struct {
	Err    error
	Status int
}

type lambda struct {
	logger   zerolog.Logger
	handlers map[string]HandleFunc
}

// NewLambda creates a new Lambda object.
func NewLambda(handlers map[string]HandleFunc) *lambda {
	logger := initLogger()

	return &lambda{
		handlers: handlers,
		logger:   logger,
	}
}

// Handle processes the incoming Lambda event
func (l *lambda) Handle(ctx context.Context, event json.RawMessage) (resp events.APIGatewayProxyResponse, err error) {
	l.logger.Debug().RawJSON("event", event).Msg("Received event")

	var base BaseRequest
	defer func() {
		if err != nil {
			l.logger.Error().Err(err).Str("action", base.Action).Msg("Prepare response error")
			err = nil
		}
	}()

	if err = serializer.UnmarshalJSON(event, &base); err != nil {
		l.logger.Error().Err(err).Str("action", base.Action).Msg("Invalid request format")
		return errResponse(http.StatusInternalServerError)
	}

	handler, ok := l.handlers[base.Action]
	if !ok {
		l.logger.Error().Err(errors.New("requested action not implemented")).Str("action", base.Action).Msg("Unknown action")
		return errResponse(http.StatusNotFound)
	}

	handlerLogger := l.logger.With().Str("action", base.Action).Logger()

	result, handleError := handler(ctx, handlerLogger, base.Data)
	if handleError != nil {
		l.logger.Error().Err(handleError.Err).Str("action", base.Action).Msg("Handle error")
		return errResponse(handleError.Status)
	}

	return okResponse(result)
}
