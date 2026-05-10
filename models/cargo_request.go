package models

import "time"

type CargoRequest struct {
	ID          string    `firestore:"id" json:"id"`
	UserID      string    `firestore:"user_id" json:"user_id"`
	UserName    string    `firestore:"user_name" json:"user_name"`
	FromCity    string    `firestore:"from_city" json:"from_city"`
	ToCity      string    `firestore:"to_city" json:"to_city"`
	Description string    `firestore:"description" json:"description"`
	WeightKg    float64   `firestore:"weight_kg" json:"weight_kg"`
	Notes       string    `firestore:"notes" json:"notes"`
	Status      string    `firestore:"status" json:"status"` // open | closed
	CreatedAt   time.Time `firestore:"created_at" json:"created_at"`
}
