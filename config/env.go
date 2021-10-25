package config

import "github.com/caarlos0/env/v6"

type Env struct {
	MainnetNode      string `env:"MAINNET_NODE"`
	TestnetNode      string `env:"TESTNET_NODE"`
	ValidatorsAppKey string `env:"VALIDATORS_APP_KEY"`
	HttpPort         uint64 `env:"HTTP_PORT" envDefault:"8080"`
}

func NewEnv() (e Env, err error) {
	err = env.Parse(&e)
	return e, err
}
