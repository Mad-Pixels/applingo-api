package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/lingo-interface/types"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func handleSubcategoriesGet(ctx context.Context, _ zerolog.Logger, _ json.RawMessage, baseParams openapi.QueryParams) (any, *api.HandleError) {
	if !api.MustGetMetaData(ctx).HasPermissions(auth.Device) {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	validSideValues := map[applingoapi.BaseSideEnum]struct{}{
		applingoapi.Front: {},
		applingoapi.Back:  {},
	}
	paramSide, err := openapi.ParseEnumParam(baseParams.GetStringPtr("side"), validSideValues)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: errors.Wrap(err, "invalid value for 'side' param")}
	}

	items := make([]applingoapi.SubcategoryItemV1, 0, len(types.AllLanguageCodes()))
	for _, code := range types.AllLanguageCodes() {
		items = append(items, applingoapi.SubcategoryItemV1{
			Code: applingoapi.BaseCountryCodeRequired(code.String()),
			Side: applingoapi.Front,
		})
	}

	response := applingoapi.CategoriesData{}
	if paramSide != nil {
		for i := range items {
			items[i].Side = *paramSide
		}
		if *paramSide == applingoapi.Front {
			response.FrontSide = items
		} else {
			response.BackSide = items
		}
	} else {
		response.FrontSide = items
		backItems := make([]applingoapi.SubcategoryItemV1, len(items))
		copy(backItems, items)
		for i := range backItems {
			backItems[i].Side = applingoapi.Back
		}
		response.BackSide = backItems
	}
	return openapi.DataResponseSubcategories(response), nil
}
