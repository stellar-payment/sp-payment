package service

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/repository"
	"github.com/stellar-payment/sp-payment/pkg/dto"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)

	// ----- Session
	AuthorizedAccessCtx(ctx context.Context, token string) (res context.Context, err error)

	// ----- Customers
	GetAllCustomer(ctx context.Context, params *dto.CustomersQueryParams) (res *dto.ListCustomerResponse, err error)
	GetCustomer(ctx context.Context, params *dto.CustomersQueryParams) (res *dto.CustomerResponse, err error)
	GetCustomerMe(ctx context.Context) (res *dto.CustomerResponse, err error)
	HandleCreateCustomer(ctx context.Context, payload *indto.EventCustomer) (err error)
	UpdateCustomer(ctx context.Context, params *dto.CustomersQueryParams, payload *dto.CustomerPayload) (err error)
	DeleteCustomer(ctx context.Context, params *dto.CustomersQueryParams) (err error)
	HandleDeleteCustomer(ctx context.Context, payload *indto.EventCustomer) (err error)

	// ----- Merchants
	GetAllMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (res *dto.ListMerchantResponse, err error)
	GetMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (res *dto.MerchantResponse, err error)
	GetMerchantMe(ctx context.Context) (res *dto.MerchantResponse, err error)
	HandleCreateMerchant(ctx context.Context, payload *indto.EventMerchant) (err error)
	UpdateMerchant(ctx context.Context, params *dto.MerchantsQueryParams, payload *dto.MerchantPayload) (err error)
	DeleteMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (err error)
	HandleDeleteMerchant(ctx context.Context, payload *indto.EventMerchant) (err error)

	// ----- Accounts
	GetAllAccount(ctx context.Context, params *dto.AccountsQueryParams) (res *dto.ListAccountResponse, err error)
	GetAccount(ctx context.Context, params *dto.AccountsQueryParams) (res *dto.AccountResponse, err error)
	GetAccountMe(ctx context.Context) (res *dto.AccountResponse, err error)
	CreateAccount(ctx context.Context, payload *dto.AccountPayload) (err error)
	UpdateAccount(ctx context.Context, params *dto.AccountsQueryParams, payload *dto.AccountPayload) (err error)
	DeleteAccount(ctx context.Context, params *dto.AccountsQueryParams) (err error)

	// ----- Transactions
	GetAllTransaction(ctx context.Context, params *dto.TransactionsQueryParams) (res *dto.ListTransactionResponse, err error)
	GetTransaction(ctx context.Context, params *dto.TransactionsQueryParams) (res *dto.TransactionResponse, err error)
	CreateTransactionP2P(ctx context.Context, payload *dto.TransactionPayload) (err error)
	CreateTransactionP2B(ctx context.Context, payload *dto.TransactionPayload) (err error)
	CreateTransactionSystem(ctx context.Context, payload *dto.TransactionPayload) (err error)
	UpdateTransaction(ctx context.Context, params *dto.TransactionsQueryParams, payload *dto.TransactionPayload) (err error)
	DeleteTransaction(ctx context.Context, params *dto.TransactionsQueryParams) (err error)

	// ----- Settlements
	GetAllSettlement(ctx context.Context, params *dto.SettlementsQueryParams) (res *dto.ListSettlementResponse, err error)
	GetSettlement(ctx context.Context, params *dto.SettlementsQueryParams) (res *dto.SettlementResponse, err error)

	// ----- Beneficiaries
	GetAllBeneficiary(ctx context.Context, params *dto.BeneficiariesQueryParams) (res *dto.ListBeneficiaryResponse, err error)
	GetBeneficiary(ctx context.Context, params *dto.BeneficiariesQueryParams) (res *dto.BeneficiaryResponse, err error)
	GetBeneficiaryPreview(ctx context.Context, params *dto.BeneficiariesQueryParams) (res float64, err error)
	CreateBeneficiary(ctx context.Context, params *dto.BeneficiariesQueryParams) (err error)
}

type service struct {
	conf       *serviceConfig
	redis      *redis.Client
	repository repository.Repository
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Repository repository.Repository
	Redis      *redis.Client
}

func NewService(params *NewServiceParams) Service {
	return &service{
		conf:       &serviceConfig{},
		repository: params.Repository,
		redis:      params.Redis,
	}
}
