package repositories

import (
	"github.com/Elfsilon/car_booking/internal/bookings/core/database"
)

type Tariffs struct {
	db *database.Database
}

func NewTariffs(db *database.Database) *Tariffs {
	return &Tariffs{db}
}

func (b *Tariffs) GetPriceByName(name string) (float64, error) {
	query := `SELECT price FROM tariffs WHERE name = $1`
	var price float64

	row := b.db.Instance().QueryRow(query, name)
	if err := row.Scan(&price); err != nil {
		return 0, err
	}

	return price, nil
}
