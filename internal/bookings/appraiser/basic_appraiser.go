package appraiser

import (
	"github.com/Elfsilon/car_booking/internal/bookings/core/utils"
)

type step struct {
	days     int
	discount float64
}

type BasicAppraiser struct {
	base  float64
	steps []step
}

func NewBasicAppraiser(base float64) RentAppraiser {
	return &BasicAppraiser{
		base: base,
		steps: []step{
			{4, 0},
			{5, 5},
			{8, 10},
			{12, 15},
		},
	}
}

func (a *BasicAppraiser) Appraise(days int) float64 {
	sum := 0.0
	for _, s := range a.steps {
		dayPrice, _ := utils.SubstractDiscount(a.base, s.discount)
		if days <= s.days {
			return sum + float64(days)*dayPrice
		}
		days -= s.days
		sum += float64(s.days) * dayPrice
	}

	return sum
}
