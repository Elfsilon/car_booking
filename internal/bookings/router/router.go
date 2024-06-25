package router

import (
	"net/http"

	"github.com/Elfsilon/car_booking/internal/bookings/controllers"
	"github.com/Elfsilon/car_booking/internal/bookings/core/database"
	"github.com/Elfsilon/car_booking/internal/bookings/repositories"
	"github.com/Elfsilon/car_booking/internal/bookings/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup(db *database.Database) *chi.Mux {
	bookingsRepo := repositories.NewBookings()
	tariffsRepo := repositories.NewTariffs(db)

	carsService := services.NewMockCars()
	bookingsService := services.NewBookings(carsService, bookingsRepo)
	tariffsService := services.NewTariffs(tariffsRepo)

	ctr := controllers.NewBookingController(tariffsService, bookingsService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	r.Get("/status", ctr.GetCarStatus)
	r.Get("/appraise", ctr.AppraisePeriod)

	return r
}
