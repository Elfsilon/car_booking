package bookings

import (
	"log"
	"net/http"

	"github.com/Elfsilon/car_booking/internal/bookings/router"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	r := router.Setup()

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
