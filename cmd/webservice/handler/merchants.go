package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/internal/util/echttputil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

type GetMerchantsHandler func(context.Context, *dto.MerchantsQueryParams) (*dto.ListMerchantResponse, error)

func HandleGetMerchants(handler GetMerchantsHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.MerchantsQueryParams{}
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

type GetMerchantByIDHandler func(context.Context, *dto.MerchantsQueryParams) (*dto.MerchantResponse, error)

func HandleGetMerchantByID(handler GetMerchantByIDHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.MerchantsQueryParams{}
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

type GetMerchantMeHandler func(context.Context) (*dto.MerchantResponse, error)

func HandleGetMerchantMe(handler GetMerchantMeHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := handler(c.Request().Context())
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}

type CreateMerchantHandler func(context.Context, *dto.MerchantPayload) error

func HandleCreateMerchants(handler CreateMerchantHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &dto.MerchantPayload{}
		if err := c.Bind(payload); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		err := handler(c.Request().Context(), payload)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

type UpdateMerchantHandler func(context.Context, *dto.MerchantsQueryParams, *dto.MerchantPayload) error

func HandleUpdateMerchants(handler UpdateMerchantHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.MerchantsQueryParams{
			MerchantID: c.Param("merchantID"),
		}

		payload := &dto.MerchantPayload{}
		if err := c.Bind(payload); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		err := handler(c.Request().Context(), params, payload)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}

type DeleteMerchantHandler func(context.Context, *dto.MerchantsQueryParams) error

func HandleDeleteMerchant(handler DeleteMerchantHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.MerchantsQueryParams{}
		if err := c.Bind(params); err != nil {
			return echttputil.WriteErrorResponse(c, errs.ErrBrokenUserReq)
		}

		err := handler(c.Request().Context(), params)
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, nil)
	}
}
