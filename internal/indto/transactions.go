package indto

import "time"

type TransactionParams struct {
	TransactionID uint64
	AccountID     string
	RecipientID   string

	TrxType   int64
	TrxTypes  []int64
	DateStart time.Time
	DateEnd   time.Time

	Keyword string
	Limit   uint64
	Page    uint64
}

type Transaction struct {
	ID            uint64    `db:"id" json:"id"`
	AccountID     string    `db:"account_id" json:"account_id"`
	AccountName   []byte    `db:"account_name" json:"account_name"`
	RecipientID   string    `db:"recipient_id" json:"recipient_id"`
	RecipientName []byte    `db:"recipient_name" json:"recipient_name"`
	TrxType       int64     `db:"trx_type" json:"trx_type"`
	TrxDatetime   time.Time `db:"trx_datetime" json:"trx_datetime"`
	TrxStatus     int64     `db:"trx_status" json:"trx_status"`
	TrxFee        float64   `db:"trx_fee" json:"trx_fee"`
	Nominal       float64   `db:"nominal" json:"nominal"`
	Description   string    `db:"description" json:"description"`
}
