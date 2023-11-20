package dto

type AccountsQueryParams struct {
	AccountID string `param:"accountID"`
	Keyword   string `query:"keyword"`
	Limit     uint64 `query:"limit"`
	Page      uint64 `query:"page"`
}

type AccountPayload struct {
	OwnerID     string  `json:"owner_id"`
	AccountType int64   `json:"account_type"`
	Balance     float64 `json:"balance"`
	AccountNo   string  `json:"account_no"`
	PIN         string  `json:"pin"`
}

type AccountResponse struct {
	ID          string  `json:"id"`
	OwnerID     string  `json:"owner_id"`
	OwnerName   string  `json:"owner_name"`
	AccountType int64   `json:"account_type"`
	Balance     float64 `json:"balance"`
	AccountNo   string  `json:"account_no"`
}

type ListAccountResponse struct {
	Accounts []*AccountResponse `json:"accounts"`
	Meta     ListPaginations    `json:"meta"`
}
