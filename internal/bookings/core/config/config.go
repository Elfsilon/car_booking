package config

type AppConfig struct {
	DB         *DatabaseConfig
	Server     *ServerConfig
	CarBooking *CarBookingConfig
}

type DatabaseConfig struct {
	ConnString string
}

type ServerConfig struct {
	Addr string
}

type CarBookingConfig struct {
	BookingPause int
}
