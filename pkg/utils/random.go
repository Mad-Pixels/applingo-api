package utils

import (
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
)

var ErrInvalidRange = errors.New("invalid range for random number")

// RandomInt returns a random integer between min and max (inclusive)
func RandomInt(min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("%w: min(%d) > max(%d)", ErrInvalidRange, min, max)
	}
	if min == max {
		return min, nil
	}
	return min + rand.Intn(max-min+1), nil
}

// RandomFromList takes a function that returns a list and returns a random element from that list.
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
