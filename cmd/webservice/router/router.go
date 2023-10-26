package router

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	ecMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/cmd/webservice/handler"
	"github.com/stellar-payment/sp-payment/internal/config"
	"github.com/stellar-payment/sp-payment/internal/middleware"
	"github.com/stellar-payment/sp-payment/internal/service"
)

type InitRouterParams struct {
	Logger  zerolog.Logger
	Service service.Service
	Ec      *echo.Echo
	Conf    *config.Config
}

func Init(params *InitRouterParams) {
	params.Ec.Use(
		ecMiddleware.CORS(), ecMiddleware.RequestIDWithConfig(ecMiddleware.RequestIDConfig{Generator: uuid.NewString}),
		middleware.ServiceVersioner,
		middleware.RequestBodyLogger(&params.Logger),
		middleware.RequestLogger(&params.Logger),
		middleware.HandlerLogger(&params.Logger),
	)

	plainRouter := params.Ec.Group("")
	secureRouter := params.Ec.Group("", middleware.AuthorizationMiddleware(params.Service))

	// ----- Maintenance
	plainRouter.GET(PingPath, handler.HandlePing(params.Service.Ping))

	// ----- Customers
	secureRouter.GET(customerBasepath, handler.HandleGetCustomers(params.Service.GetAllCustomer))
	secureRouter.OPTIONS(customerBasepath, handler.HandleGetCustomers(params.Service.GetAllCustomer))
	secureRouter.GET(customerIDPath, handler.HandleGetCustomerByID(params.Service.GetCustomer))
	secureRouter.OPTIONS(customerIDPath, handler.HandleGetCustomerByID(params.Service.GetCustomer))
	secureRouter.PUT(customerIDPath, handler.HandleUpdateCustomers(params.Service.UpdateCustomer))
	secureRouter.OPTIONS(customerIDPath, handler.HandleUpdateCustomers(params.Service.UpdateCustomer))
	secureRouter.DELETE(customerIDPath, handler.HandleDeleteCustomer(params.Service.DeleteCustomer))
	secureRouter.OPTIONS(customerIDPath, handler.HandleDeleteCustomer(params.Service.DeleteCustomer))
}
