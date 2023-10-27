package indto

import "database/sql"

type CustomerParams struct {
	UserID     string
	CustomerID string
	Keyword    string
	Limit      uint64
	Page       uint64
}

type Customer struct {
	ID           string       `db:"id"`
	UserID       string       `db:"user_id"`
	LegalName    []byte       `db:"legal_name"`
	Phone        []byte       `db:"phone"`
	Email        []byte       `db:"email"`
	Birthdate    []byte       `db:"birthdate"`
	Address      []byte       `db:"address"`
	PhotoProfile string       `db:"photo_profile"`
	RowHash      sql.NullByte `db:"row_hash"`
}

type EventCustomer struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	LegalName    string `json:"legal_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Birthdate    string `json:"birth_date"`
	Address      string `json:"address"`
	PhotoProfile string `json:"photo_profile"`
}
