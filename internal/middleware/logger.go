package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/config"
)

func HandlerLogger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			l := logger.With().Logger()
			l.UpdateContext(func(cl zerolog.Context) zerolog.Context {
				return cl.Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID))
			})

			c.SetRequest(c.Request().WithContext(l.WithContext(c.Request().Context())))
			return next(c)
		}

	}
}

func RequestLogger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID)).
				Str("latency", v.Latency.String()).
				Str("protocol", v.Protocol).
				Str("remoteIP", v.RemoteIP).
				Str("host", v.Host).
				Str("method", v.Method).
				Str("URI", v.URI).
				Str("route-path", v.RoutePath).
				Str("user-agent", v.UserAgent).
				Int("status", v.Status).Msg("request")
			return nil
		},
		LogLatency:   true,
		LogProtocol:  true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogRoutePath: true,
		LogUserAgent: true,
		LogStatus:    true,
	})
}

func RequestBodyLogger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: func(c echo.Context) bool {
			conf := config.Get()

			if conf.Environment == "prod" {
				// for security purpose, exempt body logger from auth endpoints
				return strings.Contains(c.Request().URL.Path, "/auth")
			}

			return false
		},

		Handler: func(c echo.Context, in []byte, out []byte) {
			logger.Info().
				Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID)).
				Any("request-header", c.Request().Header).Send()
			if string(in) != "" {
				logger.Info().
					Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID)).
					RawJSON("request-body", in).Send()
			}

			logger.Info().
				Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID)).
				Any("response-header", c.Response().Header()).Send()
			if string(out) != "" {
				if !strings.Contains(c.Response().Header().Get("Content-Type"), "application/json") {
					// if c.Response().Header().Get("Content-Type") != "application/json" {
					logger.Info().
						Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID)).
						Str("response-body", "<non-json response>").Send()
				} else {
					logger.Info().
						Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID)).
						RawJSON("response-body", out).Send()

				}

			}
		},
	})
}
