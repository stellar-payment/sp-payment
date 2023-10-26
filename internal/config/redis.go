package config

import "time"

type RedisConfig struct {
	Address    string        `json:"address"`
	Port       string        `json:"port"`
	Password   string        `json:"password"`
	DefaultExp time.Duration `json:"exp_duration"`
}
