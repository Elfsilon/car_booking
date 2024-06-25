package services

import (
	"github.com/Elfsilon/car_booking/internal/bookings/models"
	"github.com/Elfsilon/car_booking/internal/bookings/repositories"
)

type Bookings struct {
	cars *Cars
	rep  *repositories.Bookings
}

func NewBookings(cars *Cars, rep *repositories.Bookings) *Bookings {
	return &Bookings{cars, rep}
}

func (b *Bookings) GetCarStatus(carID string) ([]models.CarBooking, error) {
	if _, err := b.cars.GetCarInfo(carID); err != nil {
		return nil, err
	}

	bookings, err := b.rep.GetAllByCar(carID)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}
