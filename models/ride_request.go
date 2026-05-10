package models

import "time"

type RideRequest struct {
	ID             string    `firestore:"id" json:"id"`
	RideID         string    `firestore:"ride_id" json:"ride_id"`
	PassengerID    string    `firestore:"passenger_id" json:"passenger_id"`
	PassengerName  string    `firestore:"passenger_name" json:"passenger_name"`
	PassengerPhoto string    `firestore:"passenger_photo" json:"passenger_photo"`
	SeatsNeeded    int       `firestore:"seats_needed" json:"seats_needed"`
	Notes          string    `firestore:"notes" json:"notes"`
	Status         string    `firestore:"status" json:"status"` // pending | accepted | rejected | completed
	CreatedAt      time.Time `firestore:"created_at" json:"created_at"`
}
