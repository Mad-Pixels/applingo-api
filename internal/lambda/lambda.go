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
	return &lambda{
		handlers: handlers,
		logger:   initLogger(),
	}
}

// Handle processes Lambda events by routing them to specific handlers based on the "action" field.
// It expects a JSON event with "action" and "data" fields, where "action" determines the handler to use.
//
// Expected event format:
//
//	{
//	    "action": "actionName",
//	    "data": {
//	        // Action-specific payload
//	    }
//	}
//
// Example usage:
//
//	lambdaHandler := lambda.NewLambda(map[string]lambda.HandleFunc{
//		"createUser": func(ctx context.Context, logger zerolog.Logger, data json.RawMessage) (any, *lambda.HandleError) {
//			var user struct {
//				Name  string `json:"name"`
//				Email string `json:"email"`
//			}
//			if err := json.Unmarshal(data, &user); err != nil {
//				return nil, &lambda.HandleError{Err: err, Status: http.StatusBadRequest}
//			}
//			// User creation logic here
//			return user, nil
//		},
//		"getUser": func(ctx context.Context, logger zerolog.Logger, data json.RawMessage) (any, *lambda.HandleError) {
//			var request struct {
//				ID string `json:"id"`
//			}
//			if err := json.Unmarshal(data, &request); err != nil {
//				return nil, &lambda.HandleError{Err: err, Status: http.StatusBadRequest}
//			}
//			// User retrieval logic here
//			return map[string]string{"id": request.ID, "name": "John Doe"}, nil
//		},
//	})
//
// Invocation examples:
//   - Create user: {"action": "createUser", "data": {"name": "John Doe", "email": "john@example.com"}}
//   - Get user: {"action": "getUser", "data": {"id": "123"}}
func (l *lambda) Handle(ctx context.Context, event json.RawMessage) (resp events.APIGatewayProxyResponse, err error) {
	l.logger.Info().RawJSON("event", event).Msg("Received event")

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
	result, handleError := handler(ctx, l.logger, base.Data)
	if handleError != nil {
		l.logger.Error().Err(handleError.Err).Str("action", base.Action).Msg("Handle error")
		return errResponse(handleError.Status)
	}
	return okResponse(result)
}
