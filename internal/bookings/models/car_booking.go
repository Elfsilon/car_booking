package models

import "github.com/Elfsilon/car_booking/internal/pkg/date"

type CarBooking struct {
	ID   int       `json:"booking_id"`
	From date.Date `json:"date_from"`
	To   date.Date `json:"date_to"`
}

type ReportRecord struct {
	CarID       string `json:"car_id"`
	CarSign     string `json:"car_sign"`
	PercentLoad int    `json:"percent_load"`
}
