package model

import "time"

type Settlement struct {
	ID             uint64    `db:"id"`
	TransactionID  uint64    `db:"transaction_id"`
	MerchantID     string    `db:"merchant_id"`
	BeneficiaryID  uint64    `db:"beneficiary_id"`
	Amount         float64   `db:"amount"`
	SettlementDate time.Time `db:"settlement_date"`
}
