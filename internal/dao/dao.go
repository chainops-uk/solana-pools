package dao

import (
	"fmt"
	"github.com/dfuse-io/solana-go"
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
		CreatePoolValidatorData(pools ...*dmodels.PoolValidatorData) error
		SaveGovernance(gov ...*dmodels.Governance) error
		SaveCoin(coin ...*dmodels.Coin) error
		UpdatePoolData(*dmodels.PoolData) error
		UpdateValidators(validators ...*dmodels.Validator) error

		DeleteValidators(poolID uuid.UUID) error

		GetPool(name string) (*dmodels.Pool, error)
		GetLastPoolData(poolID uuid.UUID) (*dmodels.PoolData, error)
		GetLastEpochPoolData(PoolID uuid.UUID, currentEpoch uint64) (*dmodels.PoolData, error)
		GetCoinByID(id uuid.UUID) (pool *dmodels.Coin, err error)
		GetValidatorByVotePK(key solana.PublicKey) (*dmodels.Validator, error)
		GetValidator(validatorID string) (*dmodels.Validator, error)

		GetPoolCount(*postgres.Condition) (int64, error)
		GetCoinsCount(cond *postgres.Condition) (int64, error)
		GetGovernanceCount(cond *postgres.Condition) (int64, error)
		GetValidatorCount(condition *postgres.PoolValidatorDataCondition) (int64, error)

		GetPools(*postgres.Condition) ([]dmodels.Pool, error)
		GetCoins(cond *postgres.Condition) ([]*dmodels.Coin, error)
		GetGovernance(cond *postgres.Condition) ([]*dmodels.Governance, error)
		GetPoolStatistic(poolID uuid.UUID, aggregate postgres.Aggregate) ([]*dmodels.PoolData, error)
		GetValidators(condition *postgres.Condition) ([]*dmodels.Validator, error)
		GetPoolValidatorData(poolDataID uuid.UUID, condition *postgres.Condition) ([]*dmodels.PoolValidatorData, error)
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
