package model

type Customer struct {
	ID           string `db:"id"`
	UserID       string `db:"user_id"`
	LegalName    []byte `db:"legal_name"`
	Phone        []byte `db:"phone"`
	Email        []byte `db:"email"`
	Birthdate    []byte `db:"birthdate"`
	Address      []byte `db:"address"`
	PhotoProfile string `db:"photo_profile"`
	RowHash      []byte `db:"row_hash"`
}
