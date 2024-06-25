package repositories

import (
	"time"

	"github.com/Elfsilon/car_booking/internal/bookings/core/database"
	"github.com/Elfsilon/car_booking/internal/bookings/models"
)

type Bookings struct {
	db *database.Database
}

func NewBookings(db *database.Database) *Bookings {
	return &Bookings{db}
}

func (b *Bookings) GetAllByCar(carID string) ([]models.CarBooking, error) {
	query := `
		SELECT id, from_date, to_date 
		FROM bookings
		WHERE car_id = $1
	`

	rows, err := b.db.I().Query(query, carID)
	if err != nil {
		return nil, err
	}

	bookings := make([]models.CarBooking, 0)
	for rows.Next() {
		b := models.CarBooking{}
		if err := rows.Scan(&b.From, &b.To); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}

	return bookings, nil
}

func (b *Bookings) Book(userID, carID string, from, to time.Time) {
	// TODO: validate range
	// TODO: add record
}
