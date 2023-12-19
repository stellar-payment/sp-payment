package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/internal/util/echttputil"
	"github.com/stellar-payment/sp-payment/pkg/dto"
)

type GetAdminDashboardHandler func(context.Context) (*dto.AdminDashboard, error)

func HandleGetAdminDashboard(handler GetAdminDashboardHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := handler(c.Request().Context())
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}

type GetCustomerDashboardHandler func(context.Context) (*dto.CustomerDashboard, error)

func HandleGetCustomerDashboard(handler GetCustomerDashboardHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := handler(c.Request().Context())
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}

type GetMerchantDashboardHandler func(context.Context) (*dto.MerchantDashboard, error)

func HandleGetMerchantDashboard(handler GetMerchantDashboardHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := handler(c.Request().Context())
		if err != nil {
			return echttputil.WriteErrorResponse(c, err)
		}

		return echttputil.WriteSuccessResponse(c, res)
	}
}
