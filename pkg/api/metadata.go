package api

import (
	"context"
	"strconv"

	"github.com/Mad-Pixels/applingo-api/pkg/auth"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

type contextKey string

const metaDataKey contextKey = "metadata"

// GetMetaData retrieves MetaData from the given context.
// Returns the MetaData and true if found, otherwise false.
func GetMetaData(ctx context.Context) (MetaData, bool) {
	meta, ok := ctx.Value(metaDataKey).(MetaData)
	return meta, ok
}

// MustGetMetaData retrieves MetaData from the context.
// Panics if the metadata is not found.
func MustGetMetaData(ctx context.Context) MetaData {
	meta, ok := GetMetaData(ctx)
	if !ok {
		panic("metadata not found in context")
	}
	return meta
}

// MetaData holds authentication metadata extracted from the request context.
type MetaData struct {
	level      auth.Role // Role level associated with the request
	kind       auth.Kind // Type of authentication method used (e.g., JWT, HMAC)
	identifier string    // Unique identifier of the user or device
}

// HasPermissions checks if the role level is equal to or higher than the required level.
func (m MetaData) HasPermissions(requiredLevel auth.Role) bool {
	return m.level >= requiredLevel
}

// GetRole returns the role level of the request.
func (m MetaData) GetRole() auth.Role {
	return m.level
}

// IsDevice checks whether the metadata represents an authenticated device using HMAC.
func (m MetaData) IsDevice() bool {
	return m.kind == auth.HMAC && m.level == auth.Device
}

// IsUser checks whether the metadata represents an authenticated user using JWT.
func (m MetaData) IsUser() bool {
	return m.kind == auth.JWT && m.level != auth.Device
}

// ctxWithAuth extracts auth.Kind, auth.Role, and optional identifier from API Gateway request context
// and returns a new context with MetaData injected.
func ctxWithAuth(ctx context.Context, req events.APIGatewayProxyRequest) (context.Context, error) {
	kindStr, ok := req.RequestContext.Authorizer["kind"].(string)
	if !ok {
		return ctx, errors.New("missing 'kind' in context")
	}
	rawKind, err := strconv.Atoi(kindStr)
	if err != nil {
		return ctx, errors.Wrap(err, "invalid 'kind' format")
	}
	kind := auth.Kind(rawKind)
	if !auth.KindIsValid(kind) {
		return ctx, errors.New("invalid 'kind' in context")
	}

	roleStr, ok := req.RequestContext.Authorizer["role"].(string)
	if !ok {
		return ctx, errors.New("missing 'role' in context")
	}
	rawRole, err := strconv.Atoi(roleStr)
	if err != nil {
		return ctx, errors.Wrap(err, "invalid 'role' format")
	}
	level := auth.Role(rawRole)

	identifier := "ufo"
	if kind == auth.JWT {
		if id, ok := req.RequestContext.Authorizer["identifier"].(string); ok {
			identifier = id
		}
	}

	return context.WithValue(ctx, metaDataKey, MetaData{
		level:      level,
		kind:       kind,
		identifier: identifier,
	}), nil
}
