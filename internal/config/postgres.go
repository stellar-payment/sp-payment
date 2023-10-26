package config

type PostgresConfig struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`

	// Feature Flags
	FFIgnoreMigrations string
}
