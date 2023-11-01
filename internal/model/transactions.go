package model

import "time"

type Transaction struct {
	ID          uint64    `db:"id"`
	AccountID   string    `db:"account_id"`
	RecipientID string    `db:"recipient_id"`
	TrxType     int64     `db:"trx_type"`
	TrxDatetime time.Time `db:"trx_datetime"`
	TrxStatus   int64     `db:"trx_status"`
	TrxFee      float64   `db:"trx_fee"`
	Nominal     float64   `db:"nominal"`
	Description string    `db:"description"`
}
