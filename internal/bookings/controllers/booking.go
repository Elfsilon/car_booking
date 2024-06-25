package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Elfsilon/car_booking/internal/bookings/appraiser"
	"github.com/Elfsilon/car_booking/internal/bookings/services"
)

type AppraisePeriodResponse struct {
	Sum float64 `json:"sum"`
}

type BookingController struct {
	appraiser  appraiser.RentAppraiser
	carService services.CarsService
}

func NewBookingController(appraiser appraiser.RentAppraiser) *BookingController {
	return &BookingController{
		appraiser: appraiser,
	}
}

func (c *BookingController) GetCarStatus(w http.ResponseWriter, r *http.Request) {
	carID := r.URL.Query().Get("car_id")
	if carID == "" {
		http.Error(w, "param 'carID' is required", 400)
		return
	}

	if _, err := c.carService.GetCarInfo(carID); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// TODO: get car status
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
