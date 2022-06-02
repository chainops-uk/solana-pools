package dao

import (
	"fmt"
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	uuid "github.com/satori/go.uuid"
)

//go:generate moq -out postgres_mock.go . Postgres

type (
	DAO interface {
		Postgres
	}
	Postgres interface {
		CreatePoolValidatorData(pools ...*dmodels.PoolValidatorData) error
		SaveGovernance(gov ...*dmodels.Governance) error
		SaveCoin(coin ...*dmodels.Coin) error
		SaveDEFIs(defiData ...*dmodels.DEFI) error

		UpdatePoolData(*dmodels.PoolData) error
		UpdateValidators(validators ...*dmodels.Validator) error
		UpdateValidatorsData(data ...*dmodels.ValidatorData) error

		DeleteValidators(poolID uuid.UUID) error
		DeleteDeFis(cond *postgres.DeFiCondition) error

		GetPool(name string) (*dmodels.Pool, error)
		GetCoinByID(id uuid.UUID) (pool *dmodels.Coin, err error)
		GetValidator(validatorID string, epoch uint64) (*dmodels.ValidatorView, error)
		GetLastPoolDataWithApyForTenEpoch(poolID uuid.UUID) (*dmodels.PoolData, error)
		GetLastPoolData(PoolID uuid.UUID) (*dmodels.PoolData, error)
		GetValidatorByVotePK(key solana.PublicKey) (*dmodels.ValidatorView, error)
		GetLastEpochPoolData(PoolID uuid.UUID, currentEpoch uint64) (*dmodels.PoolData, error)
		GetDEFIs(cond *postgres.DeFiCondition) ([]*dmodels.DEFI, error)
		GetLiquidityPool(cond *postgres.Condition) (*dmodels.LiquidityPool, error)

		GetPoolCount(*postgres.Condition) (int64, error)
		GetCoinsCount(cond *postgres.CoinCondition) (int64, error)
		GetGovernanceCount(cond *postgres.GovernanceCondition) (int64, error)
		GetValidatorDataCount(condition *postgres.PoolValidatorDataCondition) (int64, error)
		GetValidatorCount(condition *postgres.ValidatorCondition, epoch uint64) (int64, error)
		GetLiquidityPoolsCount(cond *postgres.Condition) (int64, error)

		GetPools(condition *postgres.PoolCondition) ([]*dmodels.Pool, error)
		GetCoins(cond *postgres.CoinCondition) ([]*dmodels.Coin, error)
		GetLiquidityPools(cond *postgres.Condition) ([]*dmodels.LiquidityPool, error)
		GetGovernance(cond *postgres.GovernanceCondition) ([]*dmodels.Governance, error)
		GetValidators(condition *postgres.ValidatorCondition, epoch uint64) ([]*dmodels.ValidatorView, error)
		GetPoolStatistic(poolID uuid.UUID, aggregate postgres.Aggregate) ([]*dmodels.PoolData, error)
		GetPoolValidatorData(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error)
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
