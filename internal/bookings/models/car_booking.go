package models

import "github.com/Elfsilon/car_booking/internal/pkg/date"

type CarBooking struct {
	From date.Date `json:"date_from"`
	To   date.Date `json:"date_to"`
}
