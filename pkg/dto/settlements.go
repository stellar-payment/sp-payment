package dto

type SettlementsQueryParams struct {
	SettlementID uint64 `param:"settlementID"`
	MerchantID   string `param:"merchantID"`
	Keyword      string `query:"keyword"`
	Limit        uint64 `query:"limit"`
	Page         uint64 `query:"page"`
}

type SettlementPayload struct {
	TransactionID  uint64  `json:"transaction_id"`
	MerchantID     string  `json:"merchant_id"`
	BeneficiaryID  uint64  `json:"beneficiary_id"`
	Amount         float64 `json:"amount"`
	Status         int64   `json:"status"`
	SettlementDate string  `json:"settlement_date"`
}

type SettlementResponse struct {
	ID             uint64  `json:"id"`
	TransactionID  uint64  `json:"transaction_id"`
	MerchantID     string  `json:"merchant_id"`
	BeneficiaryID  uint64  `json:"beneficiary_id"`
	Amount         float64 `json:"amount"`
	Status         int64   `json:"status"`
	SettlementDate string  `json:"settlement_date"`
}

type ListSettlementResponse struct {
	Settlements []*SettlementResponse `json:"settlements"`
	Meta        ListPaginations       `json:"meta"`
}
