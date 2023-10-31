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
	HandleCreateCustomer(ctx context.Context, payload *indto.EventCustomer) (err error)
	UpdateCustomer(ctx context.Context, params *dto.CustomersQueryParams, payload *dto.CustomerPayload) (err error)
	DeleteCustomer(ctx context.Context, params *dto.CustomersQueryParams) (err error)
	HandleDeleteCustomer(ctx context.Context, payload *indto.EventCustomer) (err error)

	// ----- Merchants
	GetAllMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (res *dto.ListMerchantResponse, err error)
	GetMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (res *dto.MerchantResponse, err error)
	HandleCreateMerchant(ctx context.Context, payload *indto.EventMerchant) (err error)
	UpdateMerchant(ctx context.Context, params *dto.MerchantsQueryParams, payload *dto.MerchantPayload) (err error)
	DeleteMerchant(ctx context.Context, params *dto.MerchantsQueryParams) (err error)
	HandleDeleteMerchant(ctx context.Context, payload *indto.EventMerchant) (err error)

	// ----- Accounts
	GetAllAccount(ctx context.Context, params *dto.AccountsQueryParams) (res *dto.ListAccountResponse, err error)
	GetAccount(ctx context.Context, params *dto.AccountsQueryParams) (res *dto.AccountResponse, err error)
	CreateAccount(ctx context.Context, payload *dto.AccountPayload) (err error)
	UpdateAccount(ctx context.Context, params *dto.AccountsQueryParams, payload *dto.AccountPayload) (err error)
	DeleteAccount(ctx context.Context, params *dto.AccountsQueryParams) (err error)
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
