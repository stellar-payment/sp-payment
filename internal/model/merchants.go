package model

type Merchant struct {
	ID           string `db:"id"`
	UserID       string `db:"user_id"`
	Name         string `db:"name"`
	Phone        string `db:"phone"`
	Address      string `db:"address"`
	Email        string `db:"email"`
	PICName      string `db:"pic_name"`
	PICEmail     string `db:"pic_email"`
	PICPhone     string `db:"pic_phone"`
	PhotoProfile string `db:"photo_profile"`
}
