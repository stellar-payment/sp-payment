package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/internal/util/echttputil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

type GetSettlementsHandler func(context.Context, *dto.SettlementsQueryParams) (*dto.ListSettlementResponse, error)

func HandleGetSettlements(handler GetSettlementsHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.SettlementsQueryParams{}
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

type GetSettlementByIDHandler func(context.Context, *dto.SettlementsQueryParams) (*dto.SettlementResponse, error)

func HandleGetSettlementByID(handler GetSettlementByIDHandler) echo.HandlerFunc {
	return func(c echo.Context) error {

		params := &dto.SettlementsQueryParams{}
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
