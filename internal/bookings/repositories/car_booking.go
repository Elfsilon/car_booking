package repositories

import (
	"time"

	"github.com/Elfsilon/car_booking/internal/bookings/core/database"
	"github.com/Elfsilon/car_booking/internal/bookings/models"
	"github.com/Elfsilon/car_booking/internal/pkg/date"
)

type CarBooking struct {
	db *database.Database
}

func NewBookings(db *database.Database) *CarBooking {
	return &CarBooking{db}
}

// Returns bookings, where end_date == actual_end_date + booking_pause in days.
// That is to say, method returns all the days where forbidden to book a car
func (b *CarBooking) GetUnavailableDates(carID string, bookingPause int) ([]models.CarBooking, error) {
	query := `
		SELECT from_date, to_date + $1::integer
		FROM bookings
		WHERE car_id = $2
		ORDER BY to_date
	`

	rows, err := b.db.I().Query(query, bookingPause, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]models.CarBooking, 0)
	for rows.Next() {
		var from, to time.Time
		if err := rows.Scan(&from, &to); err != nil {
			return nil, err
		}
		b := models.CarBooking{
			From: date.Date{Time: from},
			To:   date.Date{Time: to},
		}
		bookings = append(bookings, b)
	}

	return bookings, nil
}

func (b *CarBooking) HasIntersections(carID string, from, to time.Time, bookingPause int) (bool, error) {
	query := `
		SELECT count(*)
		FROM bookings
		WHERE car_id = %1 AND (
			$1 >= from_date AND $2 <= to_date + $4::integer OR 
  		$2 >= from_date AND $3 <= to_date + $4::integer
		)
	`
	var intersectionsCount int
	row := b.db.I().QueryRow(query, carID, from, to, bookingPause)
	if err := row.Scan(&intersectionsCount); err != nil {
		return true, err
	}
	return intersectionsCount != 0, nil
}

func (b *CarBooking) Book(userID, carID string, from, to time.Time) (int, error) {
	query := `
		INSERT INTO bookings (user_id, car_id, from_date, to_date) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	row := b.db.I().QueryRow(query, userID, carID, from, to)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
