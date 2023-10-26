package dto

type MerchantsQueryParams struct {
	MerchantID string `param:"merchantID"`
	Keyword    string `query:"keyword"`
	Limit      uint64 `query:"limit"`
	Page       uint64 `query:"page"`
}

type MerchantPayload struct {
	Name         string `json:"name" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
	Email        string `json:"email" validate:"required"`
	Address      string `json:"address" validate:"required"`
	PICName      string `json:"pic_name" validate:"required"`
	PICEmail     string `json:"pic_email" validate:"required"`
	PICPhone     string `json:"pic_phone" validate:"required"`
	PhotoProfile string `json:"photo_profile" validate:"required"`
}

type MerchantResponse struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	PICName      string `json:"pic_name"`
	PICEmail     string `json:"pic_email"`
	PICPhone     string `json:"pic_phone"`
	PhotoProfile string `json:"photo_profile"`
}

type ListMerchantResponse struct {
	Merchants []*MerchantResponse `json:"merchants"`
	Meta      ListPaginations     `json:"meta"`
}
