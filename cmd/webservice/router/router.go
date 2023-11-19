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
	secureRouter.GET(customerMePath, handler.HandleGetCustomerMe(params.Service.GetCustomerMe))
	secureRouter.OPTIONS(customerMePath, handler.HandleGetCustomerMe(params.Service.GetCustomerMe))
	secureRouter.PUT(customerIDPath, handler.HandleUpdateCustomers(params.Service.UpdateCustomer))
	secureRouter.OPTIONS(customerIDPath, handler.HandleUpdateCustomers(params.Service.UpdateCustomer))
	secureRouter.DELETE(customerIDPath, handler.HandleDeleteCustomer(params.Service.DeleteCustomer))
	secureRouter.OPTIONS(customerIDPath, handler.HandleDeleteCustomer(params.Service.DeleteCustomer))

	// ----- Merchants
	secureRouter.GET(merchantBasepath, handler.HandleGetMerchants(params.Service.GetAllMerchant))
	secureRouter.OPTIONS(merchantBasepath, handler.HandleGetMerchants(params.Service.GetAllMerchant))
	secureRouter.GET(merchantIDPath, handler.HandleGetMerchantByID(params.Service.GetMerchant))
	secureRouter.OPTIONS(merchantIDPath, handler.HandleGetMerchantByID(params.Service.GetMerchant))
	secureRouter.GET(merchantMePath, handler.HandleGetMerchantMe(params.Service.GetMerchantMe))
	secureRouter.OPTIONS(merchantMePath, handler.HandleGetMerchantMe(params.Service.GetMerchantMe))
	secureRouter.PUT(merchantIDPath, handler.HandleUpdateMerchants(params.Service.UpdateMerchant))
	secureRouter.OPTIONS(merchantIDPath, handler.HandleUpdateMerchants(params.Service.UpdateMerchant))
	secureRouter.DELETE(merchantIDPath, handler.HandleDeleteMerchant(params.Service.DeleteMerchant))
	secureRouter.OPTIONS(merchantIDPath, handler.HandleDeleteMerchant(params.Service.DeleteMerchant))

	// ----- Accounts
	secureRouter.GET(accountBasepath, handler.HandleGetAccounts(params.Service.GetAllAccount))
	secureRouter.OPTIONS(accountBasepath, handler.HandleGetAccounts(params.Service.GetAllAccount))
	secureRouter.GET(accountIDPath, handler.HandleGetAccountByID(params.Service.GetAccount))
	secureRouter.OPTIONS(accountIDPath, handler.HandleGetAccountByID(params.Service.GetAccount))
	secureRouter.GET(accountMePath, handler.HandleGetAccountMe(params.Service.GetAccountMe))
	secureRouter.OPTIONS(accountMePath, handler.HandleGetAccountMe(params.Service.GetAccountMe))
	secureRouter.POST(accountBasepath, handler.HandleCreateAccount(params.Service.CreateAccount))
	secureRouter.OPTIONS(accountBasepath, handler.HandleCreateAccount(params.Service.CreateAccount))
	secureRouter.PUT(accountIDPath, handler.HandleUpdateAccounts(params.Service.UpdateAccount))
	secureRouter.OPTIONS(accountIDPath, handler.HandleUpdateAccounts(params.Service.UpdateAccount))
	secureRouter.DELETE(accountIDPath, handler.HandleDeleteAccount(params.Service.DeleteAccount))
	secureRouter.OPTIONS(accountIDPath, handler.HandleDeleteAccount(params.Service.DeleteAccount))

	// ----- Transactions
	secureRouter.GET(trxBasepath, handler.HandleGetTransactions(params.Service.GetAllTransaction))
	secureRouter.OPTIONS(trxBasepath, handler.HandleGetTransactions(params.Service.GetAllTransaction))
	secureRouter.GET(trxIDPath, handler.HandleGetTransactionByID(params.Service.GetTransaction))
	secureRouter.OPTIONS(trxIDPath, handler.HandleGetTransactionByID(params.Service.GetTransaction))
	secureRouter.POST(trxP2PPath, handler.HandleCreateTransaction(params.Service.CreateTransactionP2P))
	secureRouter.OPTIONS(trxP2PPath, handler.HandleCreateTransaction(params.Service.CreateTransactionP2P))
	secureRouter.POST(trxP2BPath, handler.HandleCreateTransaction(params.Service.CreateTransactionP2B))
	secureRouter.OPTIONS(trxP2BPath, handler.HandleCreateTransaction(params.Service.CreateTransactionP2B))
	secureRouter.POST(trxSYSPath, handler.HandleCreateTransaction(params.Service.CreateTransactionSystem))
	secureRouter.OPTIONS(trxSYSPath, handler.HandleCreateTransaction(params.Service.CreateTransactionSystem))
	secureRouter.PUT(trxIDPath, handler.HandleUpdateTransactions(params.Service.UpdateTransaction))
	secureRouter.OPTIONS(trxIDPath, handler.HandleUpdateTransactions(params.Service.UpdateTransaction))
	secureRouter.DELETE(trxIDPath, handler.HandleDeleteTransaction(params.Service.DeleteTransaction))
	secureRouter.OPTIONS(trxIDPath, handler.HandleDeleteTransaction(params.Service.DeleteTransaction))

	// ----- Settlements
	secureRouter.GET(settlementBasepath, handler.HandleGetSettlements(params.Service.GetAllSettlement))
	secureRouter.OPTIONS(settlementBasepath, handler.HandleGetSettlements(params.Service.GetAllSettlement))
	secureRouter.GET(settlementIDPath, handler.HandleGetSettlementByID(params.Service.GetSettlement))
	secureRouter.OPTIONS(settlementIDPath, handler.HandleGetSettlementByID(params.Service.GetSettlement))

	// ----- Beneficiaries
	secureRouter.GET(beneficiaryBasepath, handler.HandleGetBeneficiaries(params.Service.GetAllBeneficiary))
	secureRouter.OPTIONS(beneficiaryBasepath, handler.HandleGetBeneficiaries(params.Service.GetAllBeneficiary))
	secureRouter.GET(beneficiaryIDPath, handler.HandleGetBeneficiaryByID(params.Service.GetBeneficiary))
	secureRouter.OPTIONS(beneficiaryIDPath, handler.HandleGetBeneficiaryByID(params.Service.GetBeneficiary))
	secureRouter.GET(beneficiaryPreviewPath, handler.HandleGetBeneficiaryPreview(params.Service.GetBeneficiaryPreview))
	secureRouter.OPTIONS(beneficiaryPreviewPath, handler.HandleGetBeneficiaryPreview(params.Service.GetBeneficiaryPreview))
	secureRouter.POST(beneficiaryBasepath, handler.HandleCreateBeneficiary(params.Service.CreateBeneficiary))
	secureRouter.OPTIONS(beneficiaryBasepath, handler.HandleCreateBeneficiary(params.Service.CreateBeneficiary))
}
