package router

import (
	"github.com/Elfsilon/car_booking/internal/bookings/appraiser"
	"github.com/Elfsilon/car_booking/internal/bookings/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// TODO: fetch from db
const basicTariffPrice = 1000

func Setup() *chi.Mux {
	// TODO: consider DI
	apr := appraiser.NewBasicAppraiser(basicTariffPrice)
	ctr := controllers.NewBookingController(apr)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/appraise", ctr.AppraisePeriod)

	return r
}
