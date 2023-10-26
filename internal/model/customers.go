package model

import "time"

type Customer struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	LegalName    string    `db:"legal_name"`
	Phone        string    `db:"phone"`
	Email        string    `db:"email"`
	Birthdate    time.Time `db:"birth_date"`
	Address      string    `db:"address"`
	PhotoProfile string    `db:"photo_profile"`
}
