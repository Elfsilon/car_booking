package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Elfsilon/car_booking/internal/bookings/appraiser"
	"github.com/Elfsilon/car_booking/internal/bookings/models"
	"github.com/Elfsilon/car_booking/internal/bookings/router/constants"
	"github.com/Elfsilon/car_booking/internal/bookings/services"
	"github.com/Elfsilon/car_booking/internal/pkg/date"
	"github.com/go-chi/chi/v5"
)

type BookRequestBody struct {
	CarID string    `json:"car_id"`
	From  date.Date `json:"date_from"`
	To    date.Date `json:"date_to"`
}

var ErrRangeOverflow = "invalid period: range must be from 0 to 30 days long"
var ErrCarIDParamRequired = paramRequired("car_id")
var ErrFromParamRequired = paramRequired("from")
var ErrToParamRequired = paramRequired("to")
var ErrStartIsAfterEnd = "invalid period: 'from' must be earlier than 'to'"

func countDaysAndValidate(from, to time.Time) (int, error) {
	days := int(to.Sub(from).Hours()/24) + 1
	if days > 30 {
		return 0, errors.New(ErrRangeOverflow)
	}
	return days, nil
}

type CarStatusResponse struct {
	Bookings []models.CarBooking `json:"booked"`
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

// Returns only bookings without adding booking_pause before and after range
func (c *CarBookingController) GetBookedDates(w http.ResponseWriter, r *http.Request) {
	carID := r.URL.Query().Get("car_id")
	if carID == "" {
		http.Error(w, ErrCarIDParamRequired, 400)
		return
	}

	bookings, err := c.bookings.GetBookedDates(carID)
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
	userID := r.Header.Get(constants.UserIdHeaderKey)

	var body BookRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, deserializationError(err), 400)
		return
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

	id, err := c.bookings.Book(userID, body.CarID, body.From.Time, body.To.Time)
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

func (c *CarBookingController) Unbook(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(constants.UserIdHeaderKey)

	bookingIdParam := chi.URLParam(r, "booking_id")
	bookingID, err := strconv.Atoi(bookingIdParam)
	if err != nil {
		http.Error(w, invalidIntParam("booking_id"), 400)
		return
	}

	if err := c.bookings.Unbook(userID, bookingID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (c *CarBookingController) CreateReport(w http.ResponseWriter, r *http.Request) {
	year, month := r.URL.Query().Get("year"), r.URL.Query().Get("month")
	if year == "" {
		http.Error(w, paramRequired("year"), 400)
		return
	}
	if month == "" {
		http.Error(w, paramRequired("month"), 400)
		return
	}

	dateString := fmt.Sprintf("%04v-%02v-01", year, month)
	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	report, err := c.bookings.CreateReport(date)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := json.NewEncoder(w).Encode(report); err != nil {
		http.Error(w, serializationError(err), 500)
		return
	}
}
