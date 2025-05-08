package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/lingo-interface/types"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"

	"github.com/rs/zerolog"
)

func handleLevelsGet(ctx context.Context, _ zerolog.Logger, _ json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	if !api.MustGetMetaData(ctx).HasPermissions(auth.Device) {
		return nil, &api.HandleError{
			Status: http.StatusForbidden,
			Err:    errors.New("insufficient permissions"),
		}
	}

	var items []applingoapi.LevelItemV1
	for _, level := range types.AllLanguageLevels() {
		items = append(items, applingoapi.LevelItemV1{
			Level: level.String(),
		})
	}

	return openapi.DataResponseLevels(applingoapi.LevelsData{
		Items: items,
	}), nil
}
