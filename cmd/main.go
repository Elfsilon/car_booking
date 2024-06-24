package main

import "github.com/Elfsilon/car_booking/internal/bookings"

func main() {
	app := bookings.NewApp()
	app.Run()
}
