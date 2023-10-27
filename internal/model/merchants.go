package model

type Merchant struct {
	ID           string `db:"id"`
	UserID       string `db:"user_id"`
	Name         string `db:"name"`
	Phone        string `db:"phone"`
	Address      string `db:"address"`
	Email        string `db:"email"`
	PICName      []byte `db:"pic_name"`
	PICEmail     []byte `db:"pic_email"`
	PICPhone     []byte `db:"pic_phone"`
	PhotoProfile string `db:"photo_profile"`
	RowHash      []byte `db:"row_hash"`
}
