package component

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/config"
)

type InitRedisParams struct {
	Conf   *config.RedisConfig
	Logger zerolog.Logger
}

func InitRedis(params *InitRedisParams) (db *redis.Client, err error) {
	db = redis.NewClient(&redis.Options{
		Addr:     params.Conf.Address,
		Password: params.Conf.Password,
		DB:       0,
	})

	for i := 20; i > 0; i-- {
		_, err = db.Ping(context.TODO()).Result()
		if err == nil {
			break
		}

		params.Logger.Error().Msgf("error init db: %+v, retrying in 1 second", err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return
	}

	params.Logger.Info().Msg("redis init succesfully")
	return
}
