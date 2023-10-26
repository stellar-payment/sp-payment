package model

type Customer struct {
	ID           string `db:"id"`
	UserID       string `db:"user_id"`
	LegalName    string `db:"legal_name"`
	Phone        string `db:"phone"`
	Email        string `db:"email"`
	Birthdate    string `db:"birthdate"`
	Address      string `db:"address"`
	PhotoProfile string `db:"photo_profile"`
}
