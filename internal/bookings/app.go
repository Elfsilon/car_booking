package bookings

import (
	"net/http"

	"github.com/Elfsilon/car_booking/internal/bookings/core/config"
	"github.com/Elfsilon/car_booking/internal/bookings/core/database"
	"github.com/Elfsilon/car_booking/internal/bookings/router"
	"go.uber.org/zap"
)

type App struct {
	config *config.AppConfig
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	a.LoadConfig()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db := database.NewDatabase(logger, a.config.DB)
	closeDB, err := db.Open()
	if err != nil {
		logger.Fatal(err.Error())
	}
	if err := db.Instance().Ping(); err != nil {
		logger.Fatal(err.Error())
	}
	defer closeDB()

	server := http.Server{
		Addr:    a.config.Server.Addr,
		Handler: router.Setup(db),
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error(err.Error())
	}
}
