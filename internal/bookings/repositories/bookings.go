package repositories

import (
	"time"

	"github.com/Elfsilon/car_booking/internal/bookings/models"
	"github.com/Elfsilon/car_booking/internal/entity"
)

var mocks = []entity.Booking{
	{0, "A", "c84fde82-6679-4384-af11-7406de3d3e14", "2024-06-01", "2024-06-10"},
	{1, "B", "c84fde82-6679-4384-af11-7406de3d3e14", "2024-06-15", "2024-06-21"},
	{2, "C", "542af38b-c4b9-4d74-a7fa-e21c2a50a8cf", "2024-06-4", "2024-06-12"},
}

type Bookings struct{}

func NewBookings() *Bookings {
	return &Bookings{}
}

func (b *Bookings) GetAllByCar(carID string) ([]models.CarBooking, error) {
	bookings := make([]models.CarBooking, 0)
	for _, b := range mocks {
		if carID == b.CarID {
			bookings = append(bookings, models.CarBooking{
				From: b.From,
				To:   b.To,
			})
		}
	}
	return bookings, nil
}

func (b *Bookings) Book(userID, carID string, from, to time.Time) {
	// TODO: validate range
	// TODO: add record
}
