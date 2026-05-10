package models

import "time"

type Ride struct {
	ID             string    `firestore:"id" json:"id"`
	DriverID       string    `firestore:"driver_id" json:"driver_id"`
	DriverName     string    `firestore:"driver_name" json:"driver_name"`
	DriverPhoto    string    `firestore:"driver_photo" json:"driver_photo"`
	FromCity       string    `firestore:"from_city" json:"from_city"`
	ToCity         string    `firestore:"to_city" json:"to_city"`
	DepartureDate  string    `firestore:"departure_date" json:"departure_date"`
	DepartureTime  string    `firestore:"departure_time" json:"departure_time"`
	SeatsTotal     int       `firestore:"seats_total" json:"seats_total"`
	SeatsAvailable int       `firestore:"seats_available" json:"seats_available"`
	PricePerSeat   int       `firestore:"price_per_seat" json:"price_per_seat"`
	Currency       string    `firestore:"currency" json:"currency"`
	VehicleType    string    `firestore:"vehicle_type" json:"vehicle_type"`
	CargoAllowed   bool      `firestore:"cargo_allowed" json:"cargo_allowed"`
	Notes          string    `firestore:"notes" json:"notes"`
	Status         string    `firestore:"status" json:"status"` // active | completed | cancelled
	CreatedAt      time.Time `firestore:"created_at" json:"created_at"`
}
