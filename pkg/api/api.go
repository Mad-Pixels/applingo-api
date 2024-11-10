package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/pkg/logger"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
)

var (
	ErrActionNotSpecified = &HandleError{
		Err:    fmt.Errorf("action not specified in path parameters"),
		Status: http.StatusBadRequest,
	}
	ErrActionNotFound = &HandleError{
		Err:    fmt.Errorf("unknown action"),
		Status: http.StatusNotFound,
	}
)

// HandleFunc is the type for action handlers.
type HandleFunc func(context.Context, zerolog.Logger, json.RawMessage) (any, *HandleError)

// HandleError implements error from handlers.
type HandleError struct {
	Err    error
	Status int
}

// Error implements the error interface.
func (h *HandleError) Error() string {
	if h.Err != nil {
		return fmt.Sprintf("status %d: %v", h.Status, h.Err)
	}
	return fmt.Sprintf("status %d", h.Status)
}

// Config contains API configuration.
type Config struct {
	EnableRequestLogging bool
}

// API handles Lambda API Gateway proxy requests.
type API struct {
	cfg      Config
	log      zerolog.Logger
	handlers map[string]HandleFunc
}

// NewLambda creates a new API handler instance.
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

// Handle processes API Gateway proxy events, routing them to specific handlers based on the "action" path parameter.
func (a *API) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if a.cfg.EnableRequestLogging {
		a.logRequest(req)
	}

	action, handleErr := a.validateRequest(req)
	if handleErr != nil {
		return a.handleError(handleErr)
	}
	result, handleErr := a.executeHandler(ctx, action, req.Body)
	if handleErr != nil {
		return a.handleError(handleErr)
	}
	return a.createResponse(result)
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

// validateRequest validates the incoming request.
func (a *API) validateRequest(req events.APIGatewayProxyRequest) (string, *HandleError) {
	action := req.PathParameters["action"]
	if action == "" {
		return "", ErrActionNotSpecified
	}

	if _, ok := a.handlers[action]; !ok {
		return "", ErrActionNotFound
	}
	return action, nil
}

// executeHandler executes the appropriate handler for the action.
func (a *API) executeHandler(ctx context.Context, action, body string) (any, *HandleError) {
	handler := a.handlers[action]
	handlerLogger := a.log.With().Str("action", action).Logger()

	result, handleErr := handler(ctx, handlerLogger, json.RawMessage(body))
	if handleErr != nil {
		handleErr.Err = fmt.Errorf("action %s failed: %w", action, handleErr.Err)
	}
	return result, handleErr
}

// createResponse creates a successful API Gateway response.
func (a *API) createResponse(result any) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(result)
	if err != nil {
		return a.handleError(&HandleError{
			Err:    fmt.Errorf("failed to marshal response: %w", err),
			Status: http.StatusInternalServerError,
		})
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

// handleError creates an error API Gateway response.
func (a *API) handleError(herr *HandleError) (events.APIGatewayProxyResponse, error) {
	a.log.Error().Err(herr.Err).Int("status", herr.Status).Msg("Request failed")

	errBody, _ := json.Marshal(map[string]string{
		"error": herr.Error(),
	})
	return events.APIGatewayProxyResponse{
		StatusCode: herr.Status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(errBody),
	}, nil
}
