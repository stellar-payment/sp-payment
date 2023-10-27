package indto

import "database/sql"

type MerchantParams struct {
	UserID     string
	MerchantID string
	Keyword    string
	Limit      uint64
	Page       uint64
}

type Merchant struct {
	ID           string       `db:"id"`
	UserID       string       `db:"user_id"`
	Name         string       `db:"name"`
	Phone        string       `db:"phone"`
	Address      string       `db:"address"`
	Email        string       `db:"email"`
	PICName      []byte       `db:"pic_name"`
	PICEmail     []byte       `db:"pic_email"`
	PICPhone     []byte       `db:"pic_phone"`
	PhotoProfile string       `db:"photo_profile"`
	RowHash      sql.NullByte `db:"row_hash"`
}

type EventMerchant struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	PICName      string `json:"pic_name"`
	PICEmail     string `json:"pic_email"`
	PICPhone     string `json:"pic_phone"`
	PhotoProfile string `json:"photo_profile"`
}
