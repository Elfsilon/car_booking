package bookings

import (
	"fmt"
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
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("loading config")
	a.LoadConfig()

	logger.Info("Setting up db connection")
	db := database.NewDatabase(logger, a.config.DB)
	closeDB, err := db.Open()
	if err != nil {
		logger.Fatal(err.Error())
	}
	if err := db.Instance().Ping(); err != nil {
		logger.Fatal(err.Error())
	}
	defer closeDB()

	logger.Info("Initializing app")
	server := http.Server{
		Addr:    a.config.Server.Addr,
		Handler: router.Setup(a.config, db),
	}

	logger.Info(fmt.Sprintf("Running up server at %v", a.config.Server.Addr))
	if err := server.ListenAndServe(); err != nil {
		logger.Error(err.Error())
	}
}
