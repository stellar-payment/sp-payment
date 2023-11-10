package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/internal/util/echttputil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

type GetBeneficiariesHandler func(context.Context, *dto.BeneficiariesQueryParams) (*dto.ListBeneficiaryResponse, error)

func HandleGetBeneficiaries(handler GetBeneficiariesHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.BeneficiariesQueryParams{}
		if err := c.Bind(params); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		res, err := handler(c.Request().Context(), params)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}

type GetBeneficiaryByIDHandler func(context.Context, *dto.BeneficiariesQueryParams) (*dto.BeneficiaryResponse, error)

func HandleGetBeneficiaryByID(handler GetBeneficiaryByIDHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.BeneficiariesQueryParams{}
		if err := c.Bind(params); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		res, err := handler(c.Request().Context(), params)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}

type GetBeneficiaryPreviewHandler func(context.Context, *dto.BeneficiariesQueryParams) (float64, error)

func HandleGetBeneficiaryPreview(handler GetBeneficiaryPreviewHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.BeneficiariesQueryParams{}
		if err := c.Bind(params); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		res, err := handler(c.Request().Context(), params)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}

type CreateBeneficiaryHandler func(context.Context, *dto.BeneficiariesQueryParams) error

func HandleCreateBeneficiary(handler CreateBeneficiaryHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.BeneficiariesQueryParams{
			MerchantID: c.QueryParam("merchantID"),
		}

		err := handler(c.Request().Context(), params)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}
