package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Env struct {
	PostgresDSN        string `env:"POSTGRES_DSN"`
	MainnetNode        string `env:"MAINNET_NODE"`
	TestnetNode        string `env:"TESTNET_NODE"`
	ValidatorsAppKey   string `env:"VALIDATORS_APP_KEY"`
	HttpPort           uint64 `env:"HTTP_PORT" envDefault:"8080"`
	HttpSwaggerAddress string `env:"HTTP_SWAGGER_ADDRESS" envDefault:"localhost:8080"`
	GinMode            string `env:"GIN_MODE"`
}

func NewEnv() (e Env, err error) {
	err = godotenv.Load()
	if err != nil {
		return e, fmt.Errorf("can`t load env file: %s", err)
	}
	err = env.Parse(&e)
	if err != nil {
		return e, fmt.Errorf("cant` parse env file: %s", err)
	}
	return e, nil
}
