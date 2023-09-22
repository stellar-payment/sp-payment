package config

type MongoDBConfig struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"collection"`
}
