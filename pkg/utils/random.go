// Package utils provides utility functions for randomness and list handling.
package utils

import (
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
)

// ErrInvalidRange is returned when min is greater than max in a random range function.
var ErrInvalidRange = errors.New("invalid range for random number")

// RandomInt returns a random integer between min and max (inclusive).
// Returns an error if min > max.
func RandomInt(minVal, maxVal int) (int, error) {
	if minVal > maxVal {
		return 0, fmt.Errorf("%w: min(%d) > max(%d)", ErrInvalidRange, minVal, maxVal)
	}
	if minVal == maxVal {
		return minVal, nil
	}
	return minVal + rand.Intn(maxVal-minVal+1), nil
}

// RandomFromList takes a function that returns a list of strings
// and returns a random element from the list.
// Returns an error if the list is empty or retrieval fails.
func RandomFromList(listFn ListFunc) (string, error) {
	items, err := listFn()
	if err != nil {
		return "", errors.Wrap(err, "failed to get list")
	}
	if len(items) == 0 {
		return "", errors.New("empty list")
	}

	idx, err := RandomInt(0, len(items)-1)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate random index")
	}
	return items[idx], nil
}
