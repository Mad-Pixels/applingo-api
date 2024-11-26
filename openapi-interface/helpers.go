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

	DataResponseCategories = func(data applingoapi.CategoriesData) applingoapi.ResponseGetCategoriesV1 {
		return applingoapi.ResponseGetCategoriesV1{Data: data}
	}

	DataResponseDictionaries = func(data applingoapi.DictionariesData) applingoapi.ResponseGetDictionariesV1 {
		return applingoapi.ResponseGetDictionariesV1{Data: data}
	}
)
