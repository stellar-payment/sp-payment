package dto

type CustomersQueryParams struct {
	CustomerID string `param:"customerID"`
	Keyword    string `query:"keyword"`
	Limit      uint64 `query:"limit"`
	Page       uint64 `query:"page"`
}

type CustomerPayload struct {
	LegalName    string `json:"legal_name" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	Email        string `json:"email" validate:"required"`
	Birthdate    string `json:"birth_date" validate:"required"`
	Address      string `json:"address" validate:"required"`
	PhotoProfile string `json:"photo_profile"`
}

type CustomerResponse struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	LegalName    string `json:"legal_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Birthdate    string `json:"birth_date"`
	Address      string `json:"address"`
	PhotoProfile string `json:"photo_profile"`
}

type ListCustomerResponse struct {
	Customers []*CustomerResponse `json:"customers"`
	Meta      ListPaginations     `json:"meta"`
}
