package dao

import (
	"fmt"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
)

type (
	DAO interface {
		Postgres
	}
	Postgres interface {
		GetPool(name string) (pool dmodels.Pool, err error)
		GetPools() (pools []dmodels.Pool, err error)
		UpdatePool(pool dmodels.Pool) error

		GetValidators(poolID uint64) (pools []dmodels.Validator, err error)
		CreateValidators(pools []dmodels.Validator) error
		DeleteValidators(poolID uint64) error
	}
	Imp struct {
		*postgres.DB
	}
)

func NewDAO(cfg config.Env) (d DAO, err error) {
	p, err := postgres.NewDB(cfg.PostgresDSN)
	if err != nil {
		return d, fmt.Errorf("postgres.NewDB: %s", err.Error())
	}
	return &Imp{
		p,
	}, nil
}
