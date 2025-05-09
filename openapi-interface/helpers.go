package openapi

import "github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"

// DataResponseSuccess is a standard success response with message "ok".
var DataResponseSuccess = applingoapi.ResponseMessage{
	Data: applingoapi.MessageData{Message: "ok"},
}

// DataResponseMessage returns a response with a custom message.
var DataResponseMessage = func(message string) applingoapi.ResponseMessage {
	return applingoapi.ResponseMessage{Data: applingoapi.MessageData{Message: message}}
}

// DataResponseUrls returns a response containing UrlsData.
var DataResponseUrls = func(data applingoapi.UrlsData) applingoapi.ResponsePostUrlsV1 {
	return applingoapi.ResponsePostUrlsV1{Data: data}
}

// DataResponseSubcategories returns a response containing CategoriesData.
var DataResponseSubcategories = func(data applingoapi.CategoriesData) applingoapi.ResponseGetSubcategoriesV1 {
	return applingoapi.ResponseGetSubcategoriesV1{Data: data}
}

// DataResponseDictionaries returns a response containing DictionariesData.
var DataResponseDictionaries = func(data applingoapi.DictionariesData) applingoapi.ResponseGetDictionariesV1 {
	return applingoapi.ResponseGetDictionariesV1{Data: data}
}

// DataResponseLevels returns a response containing LevelsData.
var DataResponseLevels = func(data applingoapi.LevelsData) applingoapi.ResponseGetLevelsV1 {
	return applingoapi.ResponseGetLevelsV1{Data: data}
}

// DataResponseProfile returns a response containing ProfileData.
var DataResponseProfile = func(data applingoapi.ProfileData) applingoapi.ResponsePatchProfileV1 {
	return applingoapi.ResponsePatchProfileV1{Data: data}
}
