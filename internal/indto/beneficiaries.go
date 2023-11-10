package indto

import "database/sql"

type BeneficiaryParams struct {
	BeneficiaryID int64
	MerchantID    string
	Keyword       string
	Limit         uint64
	Page          uint64
}

type Beneficiary struct {
	ID             uint64       `db:"id"`
	MerchantID     string       `db:"merchant_id"`
	Amount         string       `db:"amount"`
	WithdrawalDate sql.NullTime `db:"withdrawal_date"`
	Status         int64        `db:"status"`
}
