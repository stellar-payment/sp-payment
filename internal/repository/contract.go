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
