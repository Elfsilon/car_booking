package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Elfsilon/car_booking/internal/bookings/appraiser"
	"github.com/Elfsilon/car_booking/internal/bookings/models"
	"github.com/Elfsilon/car_booking/internal/bookings/services"
	"github.com/Elfsilon/car_booking/internal/pkg/date"
)

type RentRequestBody struct {
	UserID string    `json:"user_id"`
	CarID  string    `json:"car_id"`
	From   date.Date `json:"date_from"`
	To     date.Date `json:"date_to"`
}

var ErrRangeOverflow = "invalid period: range must be from 0 to 30 days long (excluding)"

func countDaysAndValidate(from, to time.Time) (int, error) {
	days := int(to.Sub(from).Hours() / 24)
	if days > 29 {
		return 0, errors.New(ErrRangeOverflow)
	}
	return days, nil
}

type CarStatusResponse struct {
	Bookings []models.CarBooking `json:"bookings"`
}

type AppraisePeriodResponse struct {
	Sum float64 `json:"sum"`
}

type CarBookingController struct {
	tariffs  *services.Tariffs
	bookings *services.CarBooking
}

func NewBookingController(tariffs *services.Tariffs, bookings *services.CarBooking) *CarBookingController {
	return &CarBookingController{
		tariffs:  tariffs,
		bookings: bookings,
	}
}

var ErrCarIDParamRequired = paramRequired("car_id")

// Method for checking the actual available dates for booking
func (c *CarBookingController) GetUnavailableDates(w http.ResponseWriter, r *http.Request) {
	carID := r.URL.Query().Get("car_id")
	if carID == "" {
		http.Error(w, ErrCarIDParamRequired, 400)
		return
	}

	bookings, err := c.bookings.GetUnavailableDates(carID)
	if err != nil {
		if errors.Is(err, services.ErrCarNotFound) {
			http.Error(w, err.Error(), 400)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	err = json.NewEncoder(w).Encode(CarStatusResponse{bookings})
	if err != nil {
		http.Error(w, serializationError(err), 500)
		return
	}
}

var ErrFromParamRequired = paramRequired("from")
var ErrToParamRequired = paramRequired("to")
var ErrStartIsAfterEnd = "invalid period: 'from' must be earlier than 'to'"

func (c *CarBookingController) AppraisePeriod(w http.ResponseWriter, r *http.Request) {
	fromDateString := r.URL.Query().Get("from")
	if fromDateString == "" {
		http.Error(w, ErrFromParamRequired, 400)
		return
	}

	toDateString := r.URL.Query().Get("to")
	if toDateString == "" {
		http.Error(w, ErrToParamRequired, 400)
		return
	}

	fromDate, err := time.Parse(time.DateOnly, fromDateString)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	toDate, err := time.Parse(time.DateOnly, toDateString)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if fromDate.After(toDate) {
		http.Error(w, ErrStartIsAfterEnd, 400)
		return
	}

	days, err := countDaysAndValidate(fromDate, toDate)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	tariffPrice, err := c.tariffs.GetPriceByName("basic")
	if err != nil {
		http.Error(w, tariffError(err), 500)
		return
	}

	apr := appraiser.NewBasicAppraiser(tariffPrice)
	sum := apr.Appraise(days)

	err = json.NewEncoder(w).Encode(AppraisePeriodResponse{Sum: sum})
	if err != nil {
		http.Error(w, serializationError(err), 500)
		return
	}
}

func (c *CarBookingController) Book(w http.ResponseWriter, r *http.Request) {
	var body RentRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, deserializationError(err), 400)
	}

	if body.From.Time.After(body.To.Time) {
		http.Error(w, ErrStartIsAfterEnd, 400)
		return
	}

	_, err := countDaysAndValidate(body.From.Time, body.To.Time)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	id, err := c.bookings.Book(body.UserID, body.CarID, body.From.Time, body.To.Time)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(struct {
		ID int `json:"book_id"`
	}{id})
	if err != nil {
		http.Error(w, serializationError(err), 500)
		return
	}
}
