package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubstractDiscountSuccess(t *testing.T) {
	res, err := SubstractDiscount(1000, 5)
	assert.NoError(t, err)
	assert.Equal(t, 950.0, res)
}

func TestSubstractDiscountMaxPercent(t *testing.T) {
	res, err := SubstractDiscount(1000, 100)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, res)
}

func TestSubstractDiscountMinPercent(t *testing.T) {
	res, err := SubstractDiscount(1000, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, res)
}

func TestSubstractDiscountZeroValueOnError(t *testing.T) {
	res, err := SubstractDiscount(1000, -5)
	assert.Error(t, err)
	assert.Equal(t, 0.0, res)
}

func TestSubstractDiscountNegativePercent(t *testing.T) {
	_, err := SubstractDiscount(1000, -5)
	assert.Error(t, err)
	assert.Equal(t, ErrNegativePercent, err)
}

func TestSubstractDiscountTooBigPercent(t *testing.T) {
	_, err := SubstractDiscount(1000, 120)
	assert.Error(t, err)
	assert.Equal(t, ErrTooBigPercent, err)
}
