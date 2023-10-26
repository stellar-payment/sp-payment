package indto

type MerchantParams struct {
	UserID     string
	MerchantID string
	Keyword    string
	Limit      uint64
	Page       uint64
}

type Merchant struct {
	ID           string `db:"id" json:"id"`
	UserID       string `db:"user_id" json:"user_id"`
	Name         string `db:"name" json:"name"`
	Phone        string `db:"phone" json:"phone"`
	Address      string `db:"address" json:"address"`
	Email        string `db:"email" json:"email"`
	PICName      string `db:"pic_name" json:"pic_name"`
	PICEmail     string `db:"pic_email" json:"pic_email"`
	PICPhone     string `db:"pic_phone" json:"pic_phone"`
	PhotoProfile string `db:"photo_profile" json:"photo_profile"`
}
