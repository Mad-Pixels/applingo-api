package api

import (
	"context"
	"encoding/json"
	"github.com/Mad-Pixels/lingocards-api/pkg/logger"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
)

// HandleFunc is the type for action handlers.
type HandleFunc func(context.Context, zerolog.Logger, json.RawMessage) (any, *HandleError)

// HandleError implement error from handlers.
type HandleError struct {
	Err    error
	Status int
}

// Config for lambda object.
type Config struct{}

type api struct {
	cfg      Config
	log      zerolog.Logger
	handlers map[string]HandleFunc
}

// NewLambda creates a new Lambda object.
func NewLambda(cfg Config, handlers map[string]HandleFunc) *api {
	return &api{
		cfg:      cfg,
		handlers: handlers,
		log:      logger.InitLogger(),
	}
}

// Handle processes API Gateway proxy events, routing them to specific handlers based on the "action" field.
// The action is extracted from query parameters, while the handler-specific data is in the request body.
//
// Example API Gateway event:
//
//	{
//		"queryStringParameters": { "action": "example" },
//		"body": "{\"param1\":\"val1\",\"param2\":\"val2\"}"
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
//	func main() {
//		lambda.Start(lambdaHandler.Handle)
//	}
//
// Invocation examples (API Gateway request body):
//   - Create user: {"name": "John Doe", "email": "john@example.com"}
//   - Get user: {"id": "123"}
func (a *api) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (resp events.APIGatewayProxyResponse, err error) {
	a.log.Info().
		Str("path", req.Path).
		Str("httpMethod", req.HTTPMethod).
		Str("domainName", req.RequestContext.DomainName).
		Str("sourceIp", req.RequestContext.Identity.SourceIP).
		Str("userAgent", req.RequestContext.Identity.UserAgent).
		Msg("Received API Gateway event")

	action := req.PathParameters["action"]
	if action == "" {
		a.log.Error().Msg("Action not specified in path parameters")
		return errResponse(http.StatusBadRequest)
	}
	handler, ok := a.handlers[action]
	if !ok {
		a.log.Error().Str("action", action).Msg("Unknown action")
		return errResponse(http.StatusNotFound)
	}

	result, handleError := handler(ctx, a.log, json.RawMessage(req.Body))
	if handleError != nil {
		a.log.Error().Err(handleError.Err).Str("action", action).Msg("Handle error")
		return errResponse(handleError.Status)
	}
	return okResponse(result)
}
