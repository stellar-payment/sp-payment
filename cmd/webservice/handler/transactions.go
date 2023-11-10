package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/internal/util/echttputil"
	"github.com/stellar-payment/sp-payment/internal/util/structutil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

type GetTransactionsHandler func(context.Context, *dto.TransactionsQueryParams) (*dto.ListTransactionResponse, error)

func HandleGetTransactions(handler GetTransactionsHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.TransactionsQueryParams{}
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

type GetTransactionByIDHandler func(context.Context, *dto.TransactionsQueryParams) (*dto.TransactionResponse, error)

func HandleGetTransactionByID(handler GetTransactionByIDHandler) echo.HandlerFunc {
	return func(c echo.Context) error {

		params := &dto.TransactionsQueryParams{}
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

type CreateTransactionHandler func(context.Context, *dto.TransactionPayload) error

func HandleCreateTransaction(handler CreateTransactionHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &dto.TransactionPayload{}
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

type UpdateTransactionHandler func(context.Context, *dto.TransactionsQueryParams, *dto.TransactionPayload) error

func HandleUpdateTransactions(handler UpdateTransactionHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.TransactionsQueryParams{
			TransactionID: structutil.StringToUint64(c.Param("trxID")),
		}

		payload := &dto.TransactionPayload{}
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

type DeleteTransactionHandler func(context.Context, *dto.TransactionsQueryParams) error

func HandleDeleteTransaction(handler DeleteTransactionHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.TransactionsQueryParams{}
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
