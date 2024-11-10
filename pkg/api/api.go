package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/pkg/logger"
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

// Config for api object.
type Config struct {
	EnableRequestLogging bool
}

// API handles Lambda API Gateway proxy requests.
type API struct {
	cfg      Config
	log      zerolog.Logger
	handlers map[string]HandleFunc
}

// NewLambda creates a new API instance.
func NewLambda(cfg Config, handlers map[string]HandleFunc) *API {
	if handlers == nil {
		panic("handlers map cannot be nil")
	}
	return &API{
		cfg:      cfg,
		handlers: handlers,
		log:      logger.InitLogger(),
	}
}

// logRequest logs detailed request information.
func (a *API) logRequest(req events.APIGatewayProxyRequest) {
	a.log.Info().
		Str("path", req.Path).
		Str("httpMethod", req.HTTPMethod).
		Str("domainName", req.RequestContext.DomainName).
		Str("sourceIp", req.RequestContext.Identity.SourceIP).
		Str("userAgent", req.RequestContext.Identity.UserAgent).
		Msg("Received API Gateway event")
}

// Handle processes API Gateway proxy events, routing them to specific handlers based on the "action" field.
func (a *API) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if a.cfg.EnableRequestLogging {
		a.logRequest(req)
	}

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
