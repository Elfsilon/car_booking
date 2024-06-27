package services

import (
	"errors"
	"time"

	"github.com/Elfsilon/car_booking/internal/bookings/core/config"
	"github.com/Elfsilon/car_booking/internal/bookings/models"
	"github.com/Elfsilon/car_booking/internal/bookings/repositories"
)

var ErrStartOrEndAtWeekends = errors.New("booking must not starts and ends at weekends")
var ErrRangeIntersects = errors.New("picked range intersects existing bookings")

type CarBooking struct {
	config *config.CarBookingConfig
	cars   *Cars
	rep    *repositories.CarBooking
}

func NewBookings(config *config.CarBookingConfig, cars *Cars, rep *repositories.CarBooking) *CarBooking {
	return &CarBooking{config, cars, rep}
}

func (b *CarBooking) getUnavailableDates(carID string, bookingPause int) ([]models.CarBooking, error) {
	if _, err := b.cars.GetInfo(carID); err != nil {
		return nil, err
	}

	bookings, err := b.rep.GetUnavailableDates(carID, bookingPause)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (b *CarBooking) GetUnavailableDates(carID string) ([]models.CarBooking, error) {
	return b.getUnavailableDates(carID, b.config.BookingPause)
}

func (b *CarBooking) GetBookedDates(carID string) ([]models.CarBooking, error) {
	return b.getUnavailableDates(carID, 0)
}

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

func (b *CarBooking) Unbook(userID string, bookingID int) error {
	return b.rep.Unbook(userID, bookingID)
}

func (b *CarBooking) CreateReport(date time.Time) ([]models.ReportRecord, error) {
	report, err := b.rep.CreateReport(date)
	if err != nil {
		return nil, err
	}

	reportSlice := make([]models.ReportRecord, 0)
	for _, car := range b.cars.List() {
		rec := report[car.CarID]
		reportSlice = append(reportSlice, models.ReportRecord{
			CarID:       car.CarID,
			CarSign:     car.Sign,
			PercentLoad: rec.PercentLoad,
		})
	}
	return reportSlice, nil
}
