package openapi

import (
	"fmt"
	"strconv"
	"strings"
)

// QueryParams базовая структура для всех query параметров
type QueryParams struct {
	raw map[string]string
}

// NewQueryParams создает новый экземпляр QueryParams
func NewQueryParams(params map[string]string) QueryParams {
	if params == nil {
		params = make(map[string]string)
	}
	return QueryParams{raw: params}
}

// GetString возвращает строковый параметр или ошибку, если ключ не найден
func (q QueryParams) GetString(key string) (string, error) {
	v, ok := q.raw[key]
	if !ok {
		return "", fmt.Errorf("key '%s' not found", key)
	}
	return v, nil
}

// GetStringDefault возвращает строковый параметр или значение по умолчанию
func (q QueryParams) GetStringDefault(key, defaultValue string) string {
	if v, ok := q.raw[key]; ok {
		return v
	}
	return defaultValue
}

// GetBool возвращает булев параметр или ошибку, если ключ не найден
func (q QueryParams) GetBool(key string) (bool, error) {
	v, ok := q.raw[key]
	if !ok {
		return false, fmt.Errorf("key '%s' not found", key)
	}
	return strconv.ParseBool(v)
}

// GetBoolDefault возвращает булев параметр или значение по умолчанию
func (q QueryParams) GetBoolDefault(key string, defaultValue bool) bool {
	v, err := q.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return v
}

// GetInt возвращает целочисленный параметр или ошибку, если ключ не найден
func (q QueryParams) GetInt(key string) (int, error) {
	v, ok := q.raw[key]
	if !ok {
		return 0, fmt.Errorf("key '%s' not found", key)
	}
	return strconv.Atoi(v)
}

// GetIntDefault возвращает целочисленный параметр или значение по умолчанию
func (q QueryParams) GetIntDefault(key string, defaultValue int) int {
	v, err := q.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return v
}

// GetSlice возвращает массив значений, разделенных запятой, или ошибку, если ключ не найден
func (q QueryParams) GetSlice(key string) ([]string, error) {
	v, ok := q.raw[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found", key)
	}
	return strings.Split(v, ","), nil
}

// Has проверяет наличие параметра
func (q QueryParams) Has(key string) bool {
	_, ok := q.raw[key]
	return ok
}

// Raw возвращает исходную карту параметров
func (q QueryParams) Raw() map[string]string {
	return q.raw
}
