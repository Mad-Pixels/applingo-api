package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/pkg/logger"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
)

type HandleFunc func(context.Context, zerolog.Logger, json.RawMessage, openapi.QueryParams) (any, *HandleError)

type HandleError struct {
	Err     error
	Status  int
	Message string
}

type Config struct {
	EnableRequestLogging bool
}

type API struct {
	cfg      Config
	log      zerolog.Logger
	handlers map[string]HandleFunc
}

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

func (a *API) logRequest(req events.APIGatewayProxyRequest) {
	a.log.Info().
		Str("path", req.Path).
		Str("httpMethod", req.HTTPMethod).
		Str("domainName", req.RequestContext.DomainName).
		Str("sourceIp", req.RequestContext.Identity.SourceIP).
		Str("userAgent", req.RequestContext.Identity.UserAgent).
		Msg("Received API Gateway event")
}

func (a *API) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if a.cfg.EnableRequestLogging {
		a.logRequest(req)
	}

	operationKey := fmt.Sprintf("%s %s", req.HTTPMethod, req.Path)
	handler, ok := a.handlers[operationKey]
	if !ok {
		return toGatewayResponse(http.StatusNotFound, openapi.Response{
			Error: fmt.Sprintf("Unknown operation: %s", operationKey),
		})
	}

	queryParams := openapi.NewQueryParams(req.QueryStringParameters)
	result, handleError := handler(ctx, a.log, json.RawMessage(req.Body), queryParams)

	if handleError != nil {
		a.log.Error().
			Err(handleError.Err).
			Str("operation", operationKey).
			Msg("Handle error")

		errorMessage := handleError.Message
		if errorMessage == "" {
			errorMessage = http.StatusText(handleError.Status)
		}

		return toGatewayResponse(handleError.Status, openapi.Response{
			Error: errorMessage,
		})
	}

	if result == nil && req.HTTPMethod == http.MethodPost {
		return toGatewayResponse(http.StatusOK, openapi.DefaultSuccessResponse)
	}

	return toGatewayResponse(http.StatusOK, openapi.Response{
		Data: result,
	})
}

func toGatewayResponse(statusCode int, body openapi.Response) (events.APIGatewayProxyResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonBody),
	}, nil
}
