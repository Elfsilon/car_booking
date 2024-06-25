package database

import (
	"database/sql"

	"github.com/Elfsilon/car_booking/internal/bookings/core/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type CloseFn = func()

type Database struct {
	logger   *zap.Logger
	config   *config.DatabaseConfig
	instance *sql.DB
}

func NewDatabase(logger *zap.Logger, config *config.DatabaseConfig) *Database {
	return &Database{logger: logger, config: config}
}

func (db *Database) Open() (CloseFn, error) {
	conn, err := sql.Open("postgres", db.config.ConnString)
	if err != nil {
		return nil, err
	}
	db.instance = conn
	return db.close, nil
}

func (db *Database) close() {
	if err := db.instance.Close(); err != nil {
		db.logger.Error(err.Error())
	}
}

func (db *Database) Instance() *sql.DB {
	if db.instance == nil {
		panic("database's instance is nil: possibly database is not opened yet")
	}
	return db.instance
}

func (db *Database) I() *sql.DB {
	return db.Instance()
}
