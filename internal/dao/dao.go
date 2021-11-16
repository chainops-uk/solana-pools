package dao

import (
	"fmt"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	uuid "github.com/satori/go.uuid"
)

type (
	DAO interface {
		Postgres
	}
	Postgres interface {
		GetPool(name string) (pool dmodels.Pool, err error)
		GetPoolCount(*postgres.Condition) (int64, error)
		GetLastPoolData(uuid.UUID) (pool *dmodels.PoolData, err error)
		GetPools(cond *postgres.Condition) (pools []dmodels.Pool, err error)
		UpdatePoolData(pool *dmodels.PoolData) error

		GetValidators(poolDataID uuid.UUID) (pools []*dmodels.Validator, err error)
		CreateValidator(pools ...*dmodels.Validator) error
		DeleteValidators(poolID uuid.UUID) error
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
