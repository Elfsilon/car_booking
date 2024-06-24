package appraiser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAppraiserZeroDays(t *testing.T) {
	apr := NewBasicAppraiser(1000)
	res := apr.Appraise(0)
	assert.Equal(t, 0.0, res)
}

func TestBasicAppraiserFullPrice(t *testing.T) {
	apr := NewBasicAppraiser(1000)
	res := apr.Appraise(4)
	assert.Equal(t, 4000.0, res)
}

func TestBasicAppraiserWithDiscount(t *testing.T) {
	apr := NewBasicAppraiser(1000)
	res := apr.Appraise(9)
	assert.Equal(t, 8750.0, res)
}

func TestBasicAppraiserWithDoubleDiscount(t *testing.T) {
	apr := NewBasicAppraiser(1000)
	res := apr.Appraise(17)
	assert.Equal(t, 15950.0, res)
}

func TestBasicAppraiserWithMaxDiscount(t *testing.T) {
	apr := NewBasicAppraiser(1000)
	res := apr.Appraise(29)
	assert.Equal(t, 26150.0, res)
}
