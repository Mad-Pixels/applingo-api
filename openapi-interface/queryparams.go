package openapi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseEnumParam[T ~string](value *string, validValues map[T]struct{}) (*T, error) {
	if value == nil {
		return nil, nil
	}
	enumVal := T(*value)
	if _, ok := validValues[enumVal]; ok {
		return &enumVal, nil
	}
	return nil, errors.New("invalid enum value")
}

type QueryParams struct {
	raw map[string]string
}

func NewQueryParams(params map[string]string) QueryParams {
	if params == nil {
		params = make(map[string]string)
	}
	return QueryParams{raw: params}
}

func (q QueryParams) GetString(key string) (string, error) {
	v, ok := q.raw[key]
	if !ok {
		return "", fmt.Errorf("key '%s' not found", key)
	}
	return v, nil
}

func (q QueryParams) GetStringDefault(key, defaultValue string) string {
	if v, ok := q.raw[key]; ok {
		return v
	}
	return defaultValue
}

func (q QueryParams) GetStringPtr(key string) *string {
	if !q.Has(key) {
		return nil
	}
	v := q.GetStringDefault(key, "")
	return &v
}

func (q QueryParams) GetBool(key string) (bool, error) {
	v, ok := q.raw[key]
	if !ok {
		return false, fmt.Errorf("key '%s' not found", key)
	}
	return strconv.ParseBool(v)
}

func (q QueryParams) GetBoolDefault(key string, defaultValue bool) bool {
	v, err := q.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return v
}

func (q QueryParams) GetBoolPtr(key string) *bool {
	if !q.Has(key) {
		return nil
	}
	v := q.GetBoolDefault(key, false)
	return &v
}

func (q QueryParams) GetInt(key string) (int, error) {
	v, ok := q.raw[key]
	if !ok {
		return 0, fmt.Errorf("key '%s' not found", key)
	}
	return strconv.Atoi(v)
}

func (q QueryParams) GetIntDefault(key string, defaultValue int) int {
	v, err := q.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return v
}

func (q QueryParams) GetIntPtr(key string) *int {
	if !q.Has(key) {
		return nil
	}
	v := q.GetIntDefault(key, 0)
	return &v
}

func (q QueryParams) GetSlice(key string) ([]string, error) {
	v, ok := q.raw[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found", key)
	}
	return strings.Split(v, ","), nil
}

func (q QueryParams) GetSlicePtr(key string) *[]string {
	if !q.Has(key) {
		return nil
	}
	v, err := q.GetSlice(key)
	if err != nil {
		return nil
	}
	return &v
}

func (q QueryParams) Has(key string) bool {
	_, ok := q.raw[key]
	return ok
}

func (q QueryParams) Raw() map[string]string {
	return q.raw
}
