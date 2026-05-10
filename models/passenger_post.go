package models

import "time"

type PassengerPost struct {
	ID            string    `firestore:"id" json:"id"`
	PassengerID   string    `firestore:"passenger_id" json:"passenger_id"`
	PassengerName string    `firestore:"passenger_name" json:"passenger_name"`
	FromCity      string    `firestore:"from_city" json:"from_city"`
	ToCity        string    `firestore:"to_city" json:"to_city"`
	DesiredDate   string    `firestore:"desired_date" json:"desired_date"`
	SeatsNeeded   int       `firestore:"seats_needed" json:"seats_needed"`
	Notes         string    `firestore:"notes" json:"notes"`
	Status        string    `firestore:"status" json:"status"` // open | closed
	CreatedAt     time.Time `firestore:"created_at" json:"created_at"`
}
