package component

import (
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/config"
)

type InitPostgresParams struct {
	Conf   *config.PostgresConfig
	Logger zerolog.Logger
}

func InitPostgres(params *InitPostgresParams) (db *sqlx.DB, err error) {
	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		params.Conf.Username, params.Conf.Password,
		params.Conf.Address, params.Conf.DBName,
	)

	for i := 10; i > 0; i-- {
		db, err = sqlx.Connect("pgx", dataSource)
		if err == nil {
			break
		}

		params.Logger.Error().Err(err).Msgf("failed to init opening db for %s, retrying in 1 second", dataSource)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return
	}

	for i := 20; i > 0; i-- {
		err = db.Ping()
		if err == nil {
			break
		}

		params.Logger.Error().Err(err).Msgf("failed to ping db for %s, retrying in 1 second", dataSource)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return
	}

	params.Logger.Info().Msg("db init successfully")

	if params.Conf.FFIgnoreMigrations == "1" {
		return
	}

	dbMigrate, err := migrate.New("file://migrations", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		params.Conf.Username, params.Conf.Password,
		params.Conf.Address, params.Conf.DBName))
	if err != nil {
		params.Logger.Error().Err(err).Msg("failed to connect migration engine")
		return
	}

	if err = dbMigrate.Up(); err != nil && err != migrate.ErrNoChange {
		params.Logger.Error().Err(err).Msg("failed to perform migrations")
		return
	}

	rev, isDirty, err := dbMigrate.Version()
	if err != nil && err != migrate.ErrNilVersion {
		params.Logger.Error().Err(err).Msg("failed to fetch migration version")
		return
	}

	if isDirty {
		params.Logger.Warn().Msg("MariaDB migration is dirty")
	}

	if err == migrate.ErrNilVersion {
		params.Logger.Info().Msg("MariaDB Migration Version: None")
	} else {
		params.Logger.Info().Msgf("MariaDB Migration Version: %d", rev)
	}

	return
}
