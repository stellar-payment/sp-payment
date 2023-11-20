package indto

import (
	"time"
)

type SettlementParams struct {
	SettlementID  uint64
	TransactionID uint64
	BeneficiaryID uint64
	MerchantID    string
	Keyword       string
	Limit         uint64
	Page          uint64
}

type Settlement struct {
	ID             uint64    `db:"id"`
	TransactionID  uint64    `db:"transaction_id"`
	MerchantID     string    `db:"merchant_id"`
	MerchantName   string    `db:"merchant_name"`
	BeneficiaryID  uint64    `db:"beneficiary_id"`
	Amount         float64   `db:"amount"`
	Status         int64     `db:"status"`
	SettlementDate time.Time `db:"settlement_date"`
}
