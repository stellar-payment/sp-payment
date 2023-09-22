package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/internal/service"
)

func AuthorizationMiddleware(svc service.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// header := c.Request().Header

			// c.SetRequest(c.Request().Clone(ctx))
			return next(c)
		}
	}
}
