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
	HandleCreateCustomer(ctx context.Context, payload *indto.Customer) (err error)
	UpdateCustomer(ctx context.Context, params *dto.CustomersQueryParams, payload *dto.CustomerPayload) (err error)
	DeleteCustomer(ctx context.Context, params *dto.CustomersQueryParams) (err error)
	HandleDeleteCustomer(ctx context.Context, payload *indto.Customer) (err error)
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
