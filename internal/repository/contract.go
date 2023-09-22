package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
}

type repository struct {
	mariaDB *sqlx.DB
	mongoDB *mongo.Database
	redis   *redis.Client
	conf    *repositoryConfig
}

type repositoryConfig struct {
}

type NewRepositoryParams struct {
	MariaDB *sqlx.DB
	MongoDB *mongo.Database
	Redis   *redis.Client
}

func NewRepository(params *NewRepositoryParams) Repository {
	return &repository{
		conf:    &repositoryConfig{},
		mariaDB: params.MariaDB,
		mongoDB: params.MongoDB,
		redis:   params.Redis,
	}
}
