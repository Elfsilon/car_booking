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

// Returns bookings, where start_date == actual_start_date - booking_pause and
// end_date == actual_end_date + booking_pause in days. That is to say, method
// returns all the days where booking is not allowed
func (b *CarBooking) GetUnavailableDates(carID string, bookingPause int) ([]models.CarBooking, error) {
	query := `
		SELECT id, from_date - $1::integer, to_date + $1::integer
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
		var id int
		var from, to time.Time
		if err := rows.Scan(&id, &from, &to); err != nil {
			return nil, err
		}
		b := models.CarBooking{
			ID:   id,
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
		WHERE car_id = $1 AND (
			$2 >= from_date - $4::integer AND $2 <= to_date + $4::integer OR 
  		$3 >= from_date - $4::integer AND $3 <= to_date + $4::integer
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

func (b *CarBooking) Unbook(userID string, bookingID int) error {
	query := `
		DELETE FROM bookings
		WHERE user_id = $1 AND id = $2
	`
	_, err := b.db.I().Exec(query, userID, bookingID)
	return err
}

// date - important only year and month, day and time could be any
func (b *CarBooking) CreateReport(date time.Time) (map[string]models.ReportRecord, error) {
	query := `
		WITH start_end_month_dates AS (
			SELECT 
				date_trunc('month', $1::date)::date as month_start_date,
				(date_trunc('month', $1::date) + interval '1 month')::date as month_end_date
		), days_in_month AS (
			SELECT EXTRACT(epoch FROM (month_end_date::timestamp - month_start_date::timestamp)) / 86400 as days_count
			FROM start_end_month_dates
		), computed_booked_days_count AS (
			SELECT 
				car_id,
				CASE WHEN to_date >= month_end_date THEN month_end_date - 1::integer ELSE to_date END - 
					CASE WHEN from_date < month_start_date THEN month_start_date ELSE from_date END + 1::integer 
					as booked_days
			FROM bookings, start_end_month_dates
			WHERE from_date < month_end_date AND to_date >= month_start_date
		)

		SELECT car_id, ROUND(sum(booked_days) / days_count * 100) as booked_percent
		FROM computed_booked_days_count, days_in_month
		GROUP BY car_id, days_count
	`

	rows, err := b.db.I().Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	report := make(map[string]models.ReportRecord, 0)
	for rows.Next() {
		rec := models.ReportRecord{}
		if err := rows.Scan(&rec.CarID, &rec.PercentLoad); err != nil {
			return nil, err
		}
		report[rec.CarID] = rec
	}
	return report, nil
}
