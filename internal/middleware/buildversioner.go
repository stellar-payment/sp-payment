package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/stellar-payment/sp-payment/internal/config"
)

func ServiceVersioner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		conf := config.Get()

		c.Response().Header().Set("BUILD-TIME", conf.BuildTime)
		c.Response().Header().Set("BUILD-VER", conf.BuildVer)

		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}
