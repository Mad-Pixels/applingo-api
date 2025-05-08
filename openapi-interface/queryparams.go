package openapi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ParseEnumParam parses a string pointer into a strongly typed enum value.
// It returns nil if the input is nil, or an error if the value is not in validValues.
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

// QueryParams wraps raw query parameters and provides typed access methods.
type QueryParams struct {
	raw map[string]string
}

// NewQueryParams creates a new QueryParams instance from a raw map of string key-value pairs.
func NewQueryParams(params map[string]string) QueryParams {
	if params == nil {
		params = make(map[string]string)
	}
	return QueryParams{raw: params}
}

// GetString returns the value for the given key as a string or an error if not found.
func (q QueryParams) GetString(key string) (string, error) {
	v, ok := q.raw[key]
	if !ok {
		return "", fmt.Errorf("key '%s' not found", key)
	}
	return v, nil
}

// GetStringDefault returns the value for the key or the provided default if not found.
func (q QueryParams) GetStringDefault(key, defaultValue string) string {
	if v, ok := q.raw[key]; ok {
		return v
	}
	return defaultValue
}

// GetStringPtr returns the value for the key as a pointer or nil if the key is not present.
func (q QueryParams) GetStringPtr(key string) *string {
	if !q.Has(key) {
		return nil
	}
	v := q.GetStringDefault(key, "")
	return &v
}

// GetBool returns the value for the key as a boolean or an error if not found or invalid.
func (q QueryParams) GetBool(key string) (bool, error) {
	v, ok := q.raw[key]
	if !ok {
		return false, fmt.Errorf("key '%s' not found", key)
	}
	return strconv.ParseBool(v)
}

// GetBoolDefault returns the boolean value for the key or the default if not found or invalid.
func (q QueryParams) GetBoolDefault(key string, defaultValue bool) bool {
	v, err := q.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return v
}

// GetBoolPtr returns the boolean value for the key as a pointer or nil if the key is not present.
func (q QueryParams) GetBoolPtr(key string) *bool {
	if !q.Has(key) {
		return nil
	}
	v := q.GetBoolDefault(key, false)
	return &v
}

// GetInt returns the value for the key as an int or an error if not found or invalid.
func (q QueryParams) GetInt(key string) (int, error) {
	v, ok := q.raw[key]
	if !ok {
		return 0, fmt.Errorf("key '%s' not found", key)
	}
	return strconv.Atoi(v)
}

// GetIntDefault returns the int value for the key or the default if not found or invalid.
func (q QueryParams) GetIntDefault(key string, defaultValue int) int {
	v, err := q.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return v
}

// GetIntPtr returns the int value for the key as a pointer or nil if the key is not present.
func (q QueryParams) GetIntPtr(key string) *int {
	if !q.Has(key) {
		return nil
	}
	v := q.GetIntDefault(key, 0)
	return &v
}

// GetSlice returns a string slice by splitting the value for the key by commas.
func (q QueryParams) GetSlice(key string) ([]string, error) {
	v, ok := q.raw[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found", key)
	}
	return strings.Split(v, ","), nil
}

// GetSlicePtr returns a pointer to a string slice or nil if the key is not present or invalid.
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

// Has checks if the key exists in the query parameters.
func (q QueryParams) Has(key string) bool {
	_, ok := q.raw[key]
	return ok
}

// Raw returns the underlying raw map of query parameters.
func (q QueryParams) Raw() map[string]string {
	return q.raw
}
