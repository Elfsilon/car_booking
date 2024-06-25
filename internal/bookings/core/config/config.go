package config

type AppConfig struct {
	DB     *DatabaseConfig
	Server *ServerConfig
}

type DatabaseConfig struct {
	ConnString string
}

type ServerConfig struct {
	Addr string
}
