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

type CarStatusResponse struct {
	Bookings []models.CarBooking `json:"bookings"`
}

type AppraisePeriodResponse struct {
	Sum float64 `json:"sum"`
}

type BookingController struct {
	appraiser appraiser.RentAppraiser
	bookings  *services.Bookings
}

func NewBookingController(appraiser appraiser.RentAppraiser, bookings *services.Bookings) *BookingController {
	return &BookingController{
		appraiser: appraiser,
		bookings:  bookings,
	}
}

func (c *BookingController) GetCarStatus(w http.ResponseWriter, r *http.Request) {
	carID := r.URL.Query().Get("car_id")
	if carID == "" {
		http.Error(w, "param 'carID' is required", 400)
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
		message := fmt.Sprintf("unable to serialize response: %s", err)
		http.Error(w, message, 500)
		return
	}

	w.Write(responseBytes)
}

func (c *BookingController) AppraisePeriod(w http.ResponseWriter, r *http.Request) {
	fromDateString := r.URL.Query().Get("from")
	if fromDateString == "" {
		http.Error(w, "param 'from' is required", 400)
		return
	}

	toDateString := r.URL.Query().Get("to")
	if toDateString == "" {
		http.Error(w, "param 'to' is required", 400)
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
		http.Error(w, "invalid period: 'from' must be earlier than 'to'", 400)
		return
	}

	days := int(toDate.Sub(fromDate).Hours() / 24)
	sum := c.appraiser.Appraise(days)

	response := AppraisePeriodResponse{Sum: sum}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		message := fmt.Sprintf("unable to serialize response: %s", err)
		http.Error(w, message, 500)
		return
	}

	w.Write(responseBytes)
}
