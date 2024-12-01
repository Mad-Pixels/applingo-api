package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/Mad-Pixels/applingo-api/pkg/logger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type HandleFunc func(context.Context, zerolog.Logger, json.RawMessage, openapi.QueryParams, ReqCtx) (any, *HandleError)

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

func (a *API) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if a.cfg.EnableRequestLogging {
		a.logRequest(req)
	}
	opKey := fmt.Sprintf("%s %s", req.HTTPMethod, req.Path)

	permLevel, authType, user, err := getAuthData(req)
	if err != nil {
		if a.cfg.EnableRequestLogging {
			a.logError(req, opKey, err)
		}
		return gatewayResponse(
			http.StatusUnauthorized,
			openapi.DataResponseMessage(http.StatusText(http.StatusUnauthorized)),
			nil,
		)
	}
	handler, ok := a.handlers[opKey]
	if !ok {
		if a.cfg.EnableRequestLogging {
			a.logError(req, opKey, errors.New("Unknown operation"))
		}
		return gatewayResponse(
			http.StatusNotFound,
			openapi.DataResponseMessage(http.StatusText(http.StatusNotFound)),
			nil,
		)
	}

	result, handleError := handler(
		ctx,
		a.log,
		json.RawMessage(req.Body),
		openapi.NewQueryParams(req.QueryStringParameters),
		ReqCtx{
			permissionLevel: permLevel,
			authType:        authType,
			user:            user,
		},
	)
	if handleError != nil {
		if a.cfg.EnableRequestLogging {
			a.logError(req, opKey, handleError.Err)
		}
		return gatewayResponse(
			handleError.Status,
			openapi.DataResponseMessage(http.StatusText(handleError.Status)),
			nil,
		)
	}

	var status int
	switch {
	case req.HTTPMethod == "POST":
		status = http.StatusCreated
	case req.HTTPMethod == "DELETE":
		status = http.StatusNoContent
	default:
		status = http.StatusOK
	}
	return gatewayResponse(status, result, nil)
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

func (a *API) logError(req events.APIGatewayProxyRequest, opKey string, err error) {
	a.log.Error().
		Str("httpMethod", req.HTTPMethod).
		Str("operationKey", opKey).
		Str("path", req.Path).
		Err(err).
		Msg("Error handling request")
}

func getAuthData(req events.APIGatewayProxyRequest) (int, string, string, error) {
	authType, ok := req.RequestContext.Authorizer["auth_type"].(string)
	if !ok {
		return 0, "", "", errors.New("missing auth_type in context")
	}
	if authType != string(auth.HMAC) && authType != string(auth.JWT) {
		return 0, "", "", fmt.Errorf("invalid auth_type: %s", authType)
	}
	permLevel, err := strconv.Atoi(req.RequestContext.Authorizer["permissions"].(string))
	if err != nil {
		return 0, "", "", errors.New("invalid permissions in context")
	}
	user := ""
	if authType == string(auth.JWT) {
		if id, ok := req.RequestContext.Authorizer["user_id"].(string); ok {
			user = id
		}
	}
	return permLevel, authType, user, nil
}
