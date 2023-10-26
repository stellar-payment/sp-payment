package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/internal/service"
	"github.com/stellar-payment/sp-payment/internal/util/echttputil"
	"github.com/stellar-payment/sp-payment/pkg/errs"
)

func AuthorizationMiddleware(svc service.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header

			var err error
			token := header.Get("Authorization")

			if token == "" {
				return echttputil.WriteErrorResponse(c, errs.ErrNoAccess)
			}

			splittedToken := strings.Split(token, " ")
			if len(splittedToken) != 2 {
				return echttputil.WriteErrorResponse(c, errs.ErrInvalidCred)
			}

			if splittedToken[0] != "Bearer" {
				return echttputil.WriteErrorResponse(c, errs.ErrInvalidCred)
			}

			accessToken := splittedToken[1]
			if err != nil {
				return echttputil.WriteErrorResponse(c, errs.ErrInvalidCred)
			}

			ctx, err := svc.AuthorizedAccessCtx(c.Request().Context(), accessToken)
			if err != nil {
				return echttputil.WriteErrorResponse(c, errs.ErrNoAccess)
			}

			c.SetRequest(c.Request().Clone(ctx))
			return next(c)
		}
	}
}
