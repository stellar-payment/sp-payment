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

	// ----- Transactions
	FindTransactions(ctx context.Context, params *indto.TransactionParams) (res []*indto.Transaction, err error)
	CountTransactions(ctx context.Context, params *indto.TransactionParams) (res int64, err error)
	FindTransaction(ctx context.Context, params *indto.TransactionParams) (res *indto.Transaction, err error)
	CreateTransactionP2P(ctx context.Context, payload *model.Transaction) (err error)
	CreateTransactionP2B(ctx context.Context, payload *model.Transaction) (err error)
	CreateTransactionSystem(ctx context.Context, payload *model.Transaction) (err error)
	UpdateTransaction(ctx context.Context, payload *model.Transaction) (err error)
	DeleteTransaction(ctx context.Context, params *indto.TransactionParams) (err error)

	// ----- Settlements
	FindSettlements(ctx context.Context, params *indto.SettlementParams) (res []*indto.Settlement, err error)
	CountSettlements(ctx context.Context, params *indto.SettlementParams) (res int64, err error)
	FindSettlement(ctx context.Context, params *indto.SettlementParams) (res *indto.Settlement, err error)
	FindPendingSettlement(ctx context.Context, params *indto.SettlementParams) (res float64, err error)

	// ----- Beneficiaries
	FindBeneficiaries(ctx context.Context, params *indto.BeneficiaryParams) (res []*indto.Beneficiary, err error)
	CountBeneficiaries(ctx context.Context, params *indto.BeneficiaryParams) (res int64, err error)
	FindBeneficiary(ctx context.Context, params *indto.BeneficiaryParams) (res *indto.Beneficiary, err error)
	CreateBeneficiary(ctx context.Context, payload *model.Beneficiary) (res *model.Beneficiary, err error)
	UpdateBeneficiary(ctx context.Context, payload *model.Beneficiary) (err error)
	DeleteBeneficiary(ctx context.Context, params *indto.BeneficiaryParams) (err error)

	// ---- Dashboard
	FindAdminDashboard(ctx context.Context) (res *indto.AdminDashboard, err error)
	FindMerchantDashboard(ctx context.Context, param *indto.MerchantDashboardParams) (res *indto.MerchantDashboard, err error)
	FindCustomerDashboard(ctx context.Context, param *indto.CustomerDashboardParams) (res *indto.CustomerDashboard, err error)
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
