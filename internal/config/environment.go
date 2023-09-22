package config

type Environment string

const (
	EnvironmentLocal = Environment("local")
	EnvironmentDev   = Environment("dev")
	EnvironmentProd  = Environment("prod")
)
