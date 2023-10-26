package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/godruoyi/go-snowflake"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName    string          `json:"serviceName"`
	ServiceAddress string          `json:"servicePort"`
	ServiceID      string          `json:"serviceID"`
	RPCAddress     string          `json:"rpcAddress"`
	TrustedService map[string]bool `json:"trustedService"`
	Environment    Environment     `json:"environment"`

	BuildVer  string
	BuildTime string
	FilePath  string
	RunSince  time.Time

	FFJsonLogger string

	AuthServiceAddr string

	PostgresConfig PostgresConfig `json:"mariaDBConfig"`
	RedisConfig    RedisConfig    `json:"redisConfig"`
}

const logTagConfig = "[Init Config]"

var config *Config

func Init(buildTime, buildVer string) {
	godotenv.Load("conf/.env")

	conf := Config{
		ServiceName:    os.Getenv("SERVICE_NAME"),
		ServiceAddress: os.Getenv("SERVICE_ADDR"),
		ServiceID:      os.Getenv("SERVICE_ID"),
		RPCAddress:     os.Getenv("GPRC_ADDR"),
		PostgresConfig: PostgresConfig{
			Address:            os.Getenv("POSTGRES_ADDRESS"),
			Username:           os.Getenv("POSTGRES_USERNAME"),
			Password:           os.Getenv("POSTGRES_PASSWORD"),
			DBName:             os.Getenv("POSTGRES_DBNAME"),
			FFIgnoreMigrations: os.Getenv("FF_MDB_IGNORE_MIGRATIONS"),
		},
		RedisConfig: RedisConfig{
			Address:    os.Getenv("REDIS_ADDRESS"),
			Port:       os.Getenv("REDIS_PORT"),
			Password:   os.Getenv("REDIS_PASSWORD"),
			DefaultExp: 48 * time.Hour,
		},
		BuildVer:        buildVer,
		BuildTime:       buildTime,
		FilePath:        os.Getenv("FILE_PATH"),
		FFJsonLogger:    os.Getenv("FF_OVERRIDE_JSON_LOGGER"),
		AuthServiceAddr: os.Getenv("AUTH_SERVICE_ADDR"),
	}

	if conf.ServiceName == "" {
		log.Fatalf("%s service name should not be empty", logTagConfig)
	}

	if conf.ServiceAddress == "" {
		log.Fatalf("%s service port should not be empty", logTagConfig)
	}

	if conf.PostgresConfig.Address == "" || conf.PostgresConfig.DBName == "" {
		log.Fatalf("%s address and db name cannot be empty", logTagConfig)
	}

	envString := os.Getenv("ENVIRONMENT")
	if envString != "dev" && envString != "prod" && envString != "local" {
		log.Fatalf("%s environment must be either local, dev or prod, found: %s", logTagConfig, envString)
	}

	conf.Environment = Environment(envString)

	conf.TrustedService = map[string]bool{conf.ServiceID: true}
	if trusted := os.Getenv("TRUSTED_SERVICES"); trusted == "" {
		conf.TrustedService["STELLAR_HENTAI"] = true
	} else {
		for _, svc := range strings.Split(trusted, ",") {
			if _, ok := conf.TrustedService[svc]; !ok {
				conf.TrustedService[svc] = true
			}
		}
	}

	snowflake.SetMachineID(snowflake.PrivateIPToMachineID())
	snowflake.SetStartTime(time.Date(2020, 8, 1, 0, 0, 0, 0, time.UTC))

	conf.RunSince = time.Now()
	config = &conf
}

func Get() (conf *Config) {
	conf = config
	return
}
