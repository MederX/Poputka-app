package models

import "time"

type User struct {
	ID           string    `firestore:"id" json:"id"`
	TelegramID   int64     `firestore:"telegram_id" json:"telegram_id"`
	Username     string    `firestore:"username" json:"username"`
	FirstName    string    `firestore:"first_name" json:"first_name"`
	LastName     string    `firestore:"last_name" json:"last_name"`
	ProfilePhoto string    `firestore:"profile_photo" json:"profile_photo"`
	PhoneNumber  string    `firestore:"phone_number" json:"phone_number"`
	Role         string    `firestore:"role" json:"role"` // driver | passenger
	CreatedAt    time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at" json:"updated_at"`
}
