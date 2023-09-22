package config

import (
	"log"
	"os"
	"strings"
	"time"

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

	MariaDBConfig MariaDBConfig `json:"mariaDBConfig"`
	RedisConfig   RedisConfig   `json:"redisConfig"`
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
		MariaDBConfig: MariaDBConfig{
			Address:            os.Getenv("MARIADB_ADDRESS"),
			Username:           os.Getenv("MARIADB_USERNAME"),
			Password:           os.Getenv("MARIADB_PASSWORD"),
			DBName:             os.Getenv("MARIADB_DBNAME"),
			FFIgnoreMigrations: os.Getenv("FF_MDB_IGNORE_MIGRATIONS"),
		},
		RedisConfig: RedisConfig{
			Address:  os.Getenv("REDIS_ADDRESS"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		BuildVer:  buildVer,
		BuildTime: buildTime,
		FilePath:  os.Getenv("FILE_PATH"),
	}

	if conf.ServiceName == "" {
		log.Fatalf("%s service name should not be empty", logTagConfig)
	}

	if conf.ServiceAddress == "" {
		log.Fatalf("%s service port should not be empty", logTagConfig)
	}

	if conf.MariaDBConfig.Address == "" || conf.MariaDBConfig.DBName == "" {
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

	conf.RunSince = time.Now()
	config = &conf
}

func Get() (conf *Config) {
	conf = config
	return
}
