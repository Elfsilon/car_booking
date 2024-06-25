package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Elfsilon/car_booking/internal/bookings/appraiser"
	"github.com/Elfsilon/car_booking/internal/bookings/models"
	"github.com/Elfsilon/car_booking/internal/bookings/services"
)

var ErrCarIDParamRequired = paramRequired("car_id")
var ErrFromParamRequired = paramRequired("from")
var ErrToParamRequired = paramRequired("to")
var ErrStartIsAfterEnd = "invalid period: 'from' must be earlier than 'to'"
var ErrRangeOverflow = "invalid period: range must be from 0 to 30 days long (excluding)"

func paramRequired(name string) string {
	return fmt.Sprintf("param '%v' is required", name)
}

func serializationError(err error) string {
	return fmt.Sprintf("param '%v' is required", err)
}

func tariffError(err error) string {
	return fmt.Sprintf("unable to get basic tariff: %s", err)
}

type CarStatusResponse struct {
	Bookings []models.CarBooking `json:"bookings"`
}

type AppraisePeriodResponse struct {
	Sum float64 `json:"sum"`
}

type BookingController struct {
	tariffs  *services.Tariffs
	bookings *services.Bookings
}

func NewBookingController(tariffs *services.Tariffs, bookings *services.Bookings) *BookingController {
	return &BookingController{
		tariffs:  tariffs,
		bookings: bookings,
	}
}

func (c *BookingController) GetCarStatus(w http.ResponseWriter, r *http.Request) {
	carID := r.URL.Query().Get("car_id")
	if carID == "" {
		http.Error(w, ErrCarIDParamRequired, 400)
		return
	}

	bookings, err := c.bookings.GetCarStatus(carID)
	if err != nil {
		if errors.Is(err, services.ErrCarNotFound) {
			http.Error(w, err.Error(), 400)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	responseBytes, err := json.Marshal(CarStatusResponse{bookings})
	if err != nil {
		http.Error(w, serializationError(err), 500)
		return
	}

	w.Write(responseBytes)
}

func (c *BookingController) AppraisePeriod(w http.ResponseWriter, r *http.Request) {
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

	days := int(toDate.Sub(fromDate).Hours() / 24)
	if days > 29 {
		http.Error(w, ErrRangeOverflow, 400)
		return
	}

	tariffPrice, err := c.tariffs.GetPriceByName("basic")
	if err != nil {
		http.Error(w, tariffError(err), 500)
		return
	}

	apr := appraiser.NewBasicAppraiser(tariffPrice)
	sum := apr.Appraise(days)

	response := AppraisePeriodResponse{Sum: sum}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, serializationError(err), 500)
		return
	}

	w.Write(responseBytes)
}
