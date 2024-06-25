package services

import (
	"errors"
	"time"

	"github.com/Elfsilon/car_booking/internal/bookings/core/config"
	"github.com/Elfsilon/car_booking/internal/bookings/models"
	"github.com/Elfsilon/car_booking/internal/bookings/repositories"
)

type CarBooking struct {
	config *config.CarBookingConfig
	cars   *Cars
	rep    *repositories.CarBooking
}

func NewBookings(config *config.CarBookingConfig, cars *Cars, rep *repositories.CarBooking) *CarBooking {
	return &CarBooking{config, cars, rep}
}

func (b *CarBooking) GetUnavailableDates(carID string) ([]models.CarBooking, error) {
	if _, err := b.cars.GetInfo(carID); err != nil {
		return nil, err
	}

	bookings, err := b.rep.GetUnavailableDates(carID, b.config.BookingPause)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

var ErrStartOrEndAtWeekends = errors.New("booking must not starting and ending at weekends")
var ErrRangeIntersects = errors.New("picked range intersects existing bookings")

func (b *CarBooking) Book(userID, carID string, from, to time.Time) (int, error) {
	startsAtWeekends := from.Weekday() == time.Sunday || from.Weekday() == time.Saturday
	endsAtWeekends := to.Weekday() == time.Sunday || to.Weekday() == time.Saturday
	if startsAtWeekends || endsAtWeekends {
		return 0, ErrStartOrEndAtWeekends
	}

	intersects, err := b.rep.HasIntersections(carID, from, to, b.config.BookingPause)
	if err != nil {
		return 0, err
	}
	if intersects {
		return 0, ErrRangeIntersects
	}

	return b.rep.Book(userID, carID, from, to)
}
