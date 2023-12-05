package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/internal/util/echttputil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

type GetAccountsHandler func(context.Context, *dto.AccountsQueryParams) (*dto.ListAccountResponse, error)

func HandleGetAccounts(handler GetAccountsHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.AccountsQueryParams{}
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

type GetAccountByIDHandler func(context.Context, *dto.AccountsQueryParams) (*dto.AccountResponse, error)

func HandleGetAccountByID(handler GetAccountByIDHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.AccountsQueryParams{}
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

type GetAccountByAccountNoHandler func(context.Context, *dto.AccountsQueryParams) (*dto.AccountResponse, error)

func HandleGetAccountByAccountNo(handler GetAccountByAccountNoHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.AccountsQueryParams{}
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

type GetAccountMeHandler func(context.Context) (*dto.AccountResponse, error)

func HandleGetAccountMe(handler GetAccountMeHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := handler(c.Request().Context())
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}

type CreateAccountHandler func(context.Context, *dto.AccountPayload) error

func HandleCreateAccount(handler CreateAccountHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &dto.AccountPayload{}
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

type UpdateAccountHandler func(context.Context, *dto.AccountsQueryParams, *dto.AccountPayload) error

func HandleUpdateAccounts(handler UpdateAccountHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.AccountsQueryParams{
			AccountID: c.Param("accountID"),
		}

		payload := &dto.AccountPayload{}
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

type DeleteAccountHandler func(context.Context, *dto.AccountsQueryParams) error

func HandleDeleteAccount(handler DeleteAccountHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.AccountsQueryParams{}
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

type AuthenticateAccountMeHandler func(context.Context, *dto.AccountPayload) error

func HandleAuthenticateAccountMe(handler AuthenticateAccountMeHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &dto.AccountPayload{}
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
