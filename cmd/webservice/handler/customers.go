package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/internal/util/echttputil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

type GetCustomersHandler func(context.Context, *dto.CustomersQueryParams) (*dto.ListCustomerResponse, error)

func HandleGetCustomers(handler GetCustomersHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.CustomersQueryParams{}
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

type GetCustomerByIDHandler func(context.Context, *dto.CustomersQueryParams) (*dto.CustomerResponse, error)

func HandleGetCustomerByID(handler GetCustomerByIDHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.CustomersQueryParams{}
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

type CreateCustomerHandler func(context.Context, *dto.CustomerPayload) error

func HandleCreateCustomers(handler CreateCustomerHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		payload := &dto.CustomerPayload{}
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

type UpdateCustomerHandler func(context.Context, *dto.CustomersQueryParams, *dto.CustomerPayload) error

func HandleUpdateCustomers(handler UpdateCustomerHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.CustomersQueryParams{
			CustomerID: c.Param("customerID"),
		}

		payload := &dto.CustomerPayload{}
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

type DeleteCustomerHandler func(context.Context, *dto.CustomersQueryParams) error

func HandleDeleteCustomer(handler DeleteCustomerHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		params := &dto.CustomersQueryParams{}
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
