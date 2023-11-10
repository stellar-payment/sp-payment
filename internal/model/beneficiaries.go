package model

import "time"

type BeneficiaryParams struct {
	UserID     string
	MerchantID string
	Keyword    string
	Limit      uint64
	Page       uint64
}

type Beneficiary struct {
	ID             uint64     `db:"id"`
	MerchantID     string     `db:"merchant_id"`
	Amount         float64    `db:"amount"`
	WithdrawalDate *time.Time `db:"withdrawal_date"`
	Status         int64      `db:"status"`
}
