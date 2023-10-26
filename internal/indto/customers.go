package indto

type CustomerParams struct {
	UserID     string
	CustomerID string
	Keyword    string
	Limit      uint64
	Page       uint64
}

type Customer struct {
	ID           string `db:"id" json:"id"`
	UserID       string `db:"user_id" json:"user_id"`
	LegalName    string `db:"legal_name" json:"legal_name"`
	Phone        string `db:"phone" json:"phone"`
	Email        string `db:"email" json:"email"`
	Birthdate    string `db:"birthdate" json:"birth_date"`
	Address      string `db:"address" json:"address"`
	PhotoProfile string `db:"photo_profile" json:"photo_profile"`
}
