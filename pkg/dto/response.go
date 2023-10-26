package dto

import "github.com/stellar-payment/sp-payment/pkg/constant"

type ErrorResponse struct {
	Status  int              `json:"-"`
	Code    constant.ErrCode `json:"code"`
	Message string           `json:"msg"`
}

type BaseResponse struct {
	Data   interface{} `json:"data"`
	Errors interface{} `json:"error"`
}

type ListPaginations struct {
	Limit     uint64 `json:"limit"`
	Page      uint64 `json:"page"`
	TotalPage uint64 `json:"total_page"`
	TotalItem uint64 `json:"total_item,omitempty"`
}
