package utils

import (
	"errors"
)

var ErrNegativePercent = errors.New("percent value must be positive")
var ErrTooBigPercent = errors.New("percent value must not be greater than 100")

func SubstractDiscount(num float64, p float64) (float64, error) {
	if p < 0 {
		return 0.0, ErrNegativePercent
	}

	if p > 100 {
		return 0.0, ErrTooBigPercent
	}

	return num * ((100 - p) / 100), nil
}
