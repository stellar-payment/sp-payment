package dto

type TransactionsQueryParams struct {
	TransactionID uint64 `param:"trxID"`
	TrxType       int64  `query:"trxType"`
	AccountID     string `query:"accountID"`
	DateStart     string `query:"dateStart"`
	DateEnd       string `query:"dateEnd"`
	Keyword       string `query:"keyword"`
	Limit         uint64 `query:"limit"`
	Page          uint64 `query:"page"`
}

type TransactionPayload struct {
	AccountID   string  `json:"account_id" validate:"required"`
	RecipientID string  `json:"recipient_id" validate:"required"`
	TrxType     int64   `json:"trx_type" validate:"required"`
	TrxDatetime string  `json:"trx_datetime" validate:"required"`
	TrxStatus   int64   `json:"trx_status" validate:"required"`
	Nominal     float64 `json:"nominal" validate:"required"`
	Description string  `json:"description" validate:"required"`
	PIN         string  `json:"pin"`
}

type TransactionResponse struct {
	ID            uint64  `json:"id"`
	AccountID     string  `json:"account_id"`
	AccountName   string  `json:"account_name"`
	RecipientID   string  `json:"recipient_id"`
	RecipientName string  `json:"recipient_name"`
	TrxType       int64   `json:"trx_type"`
	TrxDatetime   string  `json:"trx_datetime"`
	TrxStatus     int64   `json:"trx_status"`
	TrxFee        float64 `json:"trx_fee"`
	Nominal       float64 `json:"nominal"`
	Description   string  `json:"description"`
}

type ListTransactionResponse struct {
	Transactions []*TransactionResponse `json:"transactions"`
	Meta         ListPaginations        `json:"meta"`
}
