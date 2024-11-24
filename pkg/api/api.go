package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Mad-Pixels/applingo-api/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
)

// HandleFunc is the type for action handlers.
type HandleFunc func(context.Context, zerolog.Logger, json.RawMessage, map[string]string) (any, *HandleError)

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

// Handle processes API Gateway proxy events, routing them to specific handlers based on HTTP method and path.
func (a *API) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if a.cfg.EnableRequestLogging {
		a.logRequest(req)
	}

	// Определяем операцию на основе HTTP метода
	operationId := req.HTTPMethod
	if operationId == "" {
		a.log.Error().Msg("HTTP method not specified in request")
		return errResponse(http.StatusBadRequest)
	}

	operationId = strings.ToLower(operationId) // конвертируем в нижний регистр для соответствия с operationId в OpenAPI

	handler, ok := a.handlers[operationId]
	if !ok {
		a.log.Error().
			Str("operationId", operationId).
			Str("method", req.HTTPMethod).
			Str("path", req.Path).
			Msg("Unknown operation")
		return errResponse(http.StatusNotFound)
	}

	result, handleError := handler(ctx, a.log, json.RawMessage(req.Body), req.QueryStringParameters)
	if handleError != nil {
		a.log.Error().
			Err(handleError.Err).
			Str("operationId", operationId).
			Msg("Handle error")
		return errResponse(handleError.Status)
	}
	return okResponse(result)
}
