package openapi

import "github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"

var (
	DataResponseSuccess = applingoapi.ResponseMessage{
		Data: applingoapi.MessageData{Message: "ok"},
	}

	DataResponseMessage = func(message string) applingoapi.ResponseMessage {
		return applingoapi.ResponseMessage{Data: applingoapi.MessageData{Message: message}}
	}

	DataResponseUrls = func(data applingoapi.UrlsData) applingoapi.ResponsePostUrlsV1 {
		return applingoapi.ResponsePostUrlsV1{Data: data}
	}

	DataResponseSubcategories = func(data applingoapi.CategoriesData) applingoapi.ResponseGetSubcategoriesV1 {
		return applingoapi.ResponseGetSubcategoriesV1{Data: data}
	}

	DataResponseDictionaries = func(data applingoapi.DictionariesData) applingoapi.ResponseGetDictionariesV1 {
		return applingoapi.ResponseGetDictionariesV1{Data: data}
	}

	DataResponseLevels = func(data applingoapi.LevelsData) applingoapi.ResponseGetLevelsV1 {
		return applingoapi.ResponseGetLevelsV1{Data: data}
	}

	DataResponseProfile = func(data applingoapi.ProfileData) applingoapi.ResponsePatchProfileV1 {
		return applingoapi.ResponsePatchProfileV1{Data: data}
	}
)
