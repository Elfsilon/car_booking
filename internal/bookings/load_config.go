package bookings

import (
	"fmt"

	"github.com/Elfsilon/car_booking/internal/bookings/core/config"
	"github.com/Elfsilon/car_booking/internal/pkg/env"
)

func (a *App) LoadConfig() {
	conn := env.MustLoadString("CONN_STRING")

	host, port := env.MustLoadString("HOST"), env.MustLoadInt("PORT")
	addr := fmt.Sprintf("%v:%v", host, port)

	a.config = &config.AppConfig{
		DB: &config.DatabaseConfig{
			ConnString: conn,
		},
		Server: &config.ServerConfig{
			Addr: addr,
		},
	}
}
