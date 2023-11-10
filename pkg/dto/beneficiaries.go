package dto

type BeneficiariesQueryParams struct {
	BeneficiaryID int64  `param:"beneficiaryID"`
	MerchantID    string `query:"merchantID"`
	Keyword       string `query:"keyword"`
	Limit         uint64 `query:"limit"`
	Page          uint64 `query:"page"`
}

type BeneficiaryPayload struct {
	MerchantID     string `json:"merchant_id"`
	Amount         string `json:"amount"`
	WithdrawalDate string `json:"withdrawal_date"`
	Status         int64  `json:"status"`
}

type BeneficiaryResponse struct {
	ID             uint64 `json:"id"`
	MerchantID     string `json:"merchant_id"`
	Amount         string `json:"amount"`
	WithdrawalDate string `json:"withdrawal_date"`
	Status         int64  `json:"status"`
}

type ListBeneficiaryResponse struct {
	Beneficiaries []*BeneficiaryResponse `json:"beneficiaries"`
	Meta          ListPaginations        `json:"meta"`
}
