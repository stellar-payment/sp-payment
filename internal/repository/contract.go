package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/stellar-payment/sp-payment/internal/indto"
	"github.com/stellar-payment/sp-payment/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	// ----- Customers
	FindCustomers(ctx context.Context, params *indto.CustomerParams) (res []*indto.Customer, err error)
	CountCustomers(ctx context.Context, params *indto.CustomerParams) (res int64, err error)
	FindCustomer(ctx context.Context, params *indto.CustomerParams) (res *indto.Customer, err error)
	CreateCustomer(ctx context.Context, payload *model.Customer) (res *model.Customer, err error)
	UpdateCustomer(ctx context.Context, payload *model.Customer) (err error)
	DeleteCustomer(ctx context.Context, params *indto.CustomerParams) (err error)

	// ----- Merchants
	FindMerchants(ctx context.Context, params *indto.MerchantParams) (res []*indto.Merchant, err error)
	CountMerchants(ctx context.Context, params *indto.MerchantParams) (res int64, err error)
	FindMerchant(ctx context.Context, params *indto.MerchantParams) (res *indto.Merchant, err error)
	CreateMerchant(ctx context.Context, payload *model.Merchant) (res *model.Merchant, err error)
	UpdateMerchant(ctx context.Context, payload *model.Merchant) (err error)
	DeleteMerchant(ctx context.Context, params *indto.MerchantParams) (err error)

	// ----- Accounts
	FindAccounts(ctx context.Context, params *indto.AccountParams) (res []*indto.Account, err error)
	CountAccounts(ctx context.Context, params *indto.AccountParams) (res int64, err error)
	FindAccount(ctx context.Context, params *indto.AccountParams) (res *indto.Account, err error)
	CreateAccount(ctx context.Context, payload *model.Account) (res *model.Account, err error)
	UpdateAccount(ctx context.Context, payload *model.Account) (err error)
	DeleteAccount(ctx context.Context, params *indto.AccountParams) (err error)
}

type repository struct {
	db    *sqlx.DB
	redis *redis.Client
	conf  *repositoryConfig
}

type repositoryConfig struct {
}

type NewRepositoryParams struct {
	DB      *sqlx.DB
	MongoDB *mongo.Database
	Redis   *redis.Client
}

func NewRepository(params *NewRepositoryParams) Repository {
	return &repository{
		conf:  &repositoryConfig{},
		db:    params.DB,
		redis: params.Redis,
	}
}

var pgSquirrel = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
