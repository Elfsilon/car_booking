package router

import (
	"github.com/Elfsilon/car_booking/internal/bookings/controllers"
	"github.com/Elfsilon/car_booking/internal/bookings/core/config"
	"github.com/Elfsilon/car_booking/internal/bookings/core/database"
	"github.com/Elfsilon/car_booking/internal/bookings/repositories"
	appmiddleware "github.com/Elfsilon/car_booking/internal/bookings/router/middleware"
	"github.com/Elfsilon/car_booking/internal/bookings/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup(config *config.AppConfig, db *database.Database) *chi.Mux {
	bookingsRepo := repositories.NewBookings(db)
	tariffsRepo := repositories.NewTariffs(db)

	carsService := services.NewMockCars()
	bookingsService := services.NewBookings(config.CarBooking, carsService, bookingsRepo)
	tariffsService := services.NewTariffs(tariffsRepo)

	ctr := controllers.NewBookingController(tariffsService, bookingsService)

	rootRouter := chi.NewRouter()

	rootRouter.Use(middleware.Logger)
	rootRouter.Use(appmiddleware.DefaultHeadersSetter)

	rootRouter.Get("/booked", ctr.GetUnavailableDates)
	rootRouter.Get("/appraise", ctr.AppraisePeriod)

	bookRouter := chi.NewRouter()
	bookRouter.Use(appmiddleware.UserIdHeaderChecker)

	bookRouter.Post("/", ctr.Book)
	bookRouter.Delete("/{booking_id}", ctr.Unbook)
	rootRouter.Mount("/book", bookRouter)

	return rootRouter
}
