package services

import "github.com/Elfsilon/car_booking/internal/bookings/repositories"

type Tariffs struct {
	rep *repositories.Tariffs
}

func NewTariffs(rep *repositories.Tariffs) *Tariffs {
	return &Tariffs{rep}
}

func (t *Tariffs) GetPriceByName(name string) (float64, error) {
	price, err := t.rep.GetPriceByName(name)
	if err != nil {
		return 0, err
	}
	return price, nil
}
