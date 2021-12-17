package services

import (
	"errors"
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/cache"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services/smodels"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

func (s *Imp) GetPool(name string) (*smodels.PoolDetails, error) {
	pd, err := s.cache.GetPool(name)
	if err != nil && !errors.Is(err, cache.KeyWasNotFound) {
		return nil, err
	}
	if pd != nil {
		return pd, nil
	}

	dPool, err := s.dao.GetPool(name)
	if err != nil {
		return nil, fmt.Errorf("dao.GetPool: %s", err.Error())
	}
	if dPool == nil {
		return nil, fmt.Errorf("dao.GetPool(%s): %w", name, postgres.ErrorRecordNotFounded)
	}
	dLastPoolData, err := s.dao.GetLastPoolData(dPool.ID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetPoolData: %s", err.Error())
	}
	dValidators, err := s.dao.GetPoolValidatorData(dLastPoolData.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("dao.GetPoolValidatorData: %s", err.Error())
	}
	validatorsS := make([]*smodels.Validator, len(dValidators))
	validatorsD := make([]*dmodels.Validator, len(dValidators))
	for i, v := range dValidators {
		validatorsD[i], err = s.dao.GetValidator(v.ValidatorID)
		if err != nil {
			return nil, fmt.Errorf("dao.GetValidator(%s): %w", err)
		}
		validatorsS[i] = (&smodels.Validator{}).Set(v.ActiveStake, validatorsD[i])
	}

	coin, err := s.dao.GetCoinByID(dPool.CoinID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetCoinByID: %w", err)
	}
	Pool := (&smodels.Pool{}).Set(dLastPoolData, coin, dPool, validatorsD)

	pd = &smodels.PoolDetails{
		Pool: *Pool,
	}

	s.cache.SetPool(pd, time.Second*30)

	return pd, nil
}

func (s *Imp) GetPools(name string, limit uint64, offset uint64) ([]*smodels.PoolDetails, uint64, error) {
	dPools, err := s.dao.GetPools(&postgres.Condition{
		Network: postgres.MainNet,
		Name:    name,
		Pagination: postgres.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetPool: %w", err)
	}
	if len(dPools) == 0 {
		return nil, 0, nil
	}
	pools := make([]*smodels.PoolDetails, len(dPools))
	for i, v1 := range dPools {
		pools[i] = &smodels.PoolDetails{
			Pool: smodels.Pool{
				Address: v1.Address,
				Name:    v1.Name,
			},
		}

		dLastPoolData, err := s.dao.GetLastPoolData(v1.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("dao.GetLastPoolData: %w", err)
		}

		dValidators, err := s.dao.GetPoolValidatorData(dLastPoolData.ID, nil)
		if err != nil {
			return nil, 0, fmt.Errorf("dao.GetPoolValidatorData: %w", err)
		}

		validatorsD := make([]*dmodels.Validator, len(dValidators))
		for i, v2 := range dValidators {
			validatorsD[i], err = s.dao.GetValidator(v2.ValidatorID)
			if err != nil {
				return nil, 0, fmt.Errorf("dao.GetValidator: %w", err)
			}
		}

		coin, err := s.dao.GetCoinByID(v1.CoinID)
		if err != nil {
			return nil, 0, fmt.Errorf("dao.GetCoinByID: %w", err)
		}

		pools[i].Set(dLastPoolData, coin, &v1, validatorsD)
	}

	count, err := s.dao.GetPoolCount(&postgres.Condition{
		Network: postgres.MainNet,
		Name:    name,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetPoolCount: %w", err)
	}

	return pools, uint64(count), nil
}

func (s *Imp) GetPoolsCurrentStatistic() (*smodels.Statistic, error) {
	stat, err := s.cache.GetCurrentStatistic()
	if err != nil && !errors.Is(err, cache.KeyWasNotFound) {
		return nil, err
	}
	if stat != nil {
		return stat, nil
	}

	dPools, err := s.dao.GetPools(&postgres.Condition{Network: postgres.MainNet})
	if err != nil {
		return nil, fmt.Errorf("dao.GetPool: %w", err)
	}
	if len(dPools) == 0 {
		return nil, nil
	}
	stat = &smodels.Statistic{}

	once := sync.Once{}
	pools := make([]*smodels.PoolDetails, len(dPools))

	var ActiveStakeSum, UnstakeSum, SupplySum uint64
	for i, v1 := range dPools {
		pools[i] = &smodels.PoolDetails{
			Pool: smodels.Pool{
				Address: v1.Address,
				Name:    v1.Name,
			},
		}

		dLastPoolData, err := s.dao.GetLastPoolData(v1.ID)
		if err != nil {
			return nil, fmt.Errorf("dao.GetLastPoolData: %w", err)
		}

		dValidators, err := s.dao.GetPoolValidatorData(dLastPoolData.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("dao.GetValidators: %w", err)
		}

		validatorsD := make([]*dmodels.Validator, len(dValidators))
		for i, v2 := range dValidators {
			validatorsD[i], err = s.dao.GetValidator(v2.ValidatorID)
			if err != nil {
				return nil, fmt.Errorf("dao.GetValidator: %w", err)
			}
		}

		coin, err := s.dao.GetCoinByID(v1.CoinID)
		if err != nil {
			return nil, fmt.Errorf("dao.GetCoinByID: %w", err)
		}

		pools[i].Set(dLastPoolData, coin, &v1, validatorsD)

		once.Do(func() {
			stat.MINScore = pools[i].AVGScore
			stat.MAXScore = pools[i].AVGScore
		})

		if pools[i].AVGScore > stat.MAXScore {
			stat.MAXScore = pools[i].AVGScore
		}
		if pools[i].AVGScore < stat.MINScore {
			stat.MINScore = pools[i].AVGScore
		}

		ActiveStakeSum += dLastPoolData.ActiveStake
		SupplySum += dLastPoolData.TotalTokensSupply
		UnstakeSum += dLastPoolData.UnstakeLiquidity
		stat.AVGSkippedSlots = stat.AVGSkippedSlots.Add(pools[i].AVGSkippedSlots)
		stat.AVGScore += pools[i].AVGScore
		stat.Delinquent = stat.Delinquent.Add(pools[i].Delinquent)
		stat.AVGPoolsApy = stat.AVGPoolsApy.Add(pools[i].APY)
	}

	if len(dPools) > 0 {
		stat.AVGPoolsApy.Div(decimal.NewFromInt(int64(len(dPools))))
		stat.AVGSkippedSlots = stat.AVGSkippedSlots.Div(decimal.NewFromInt(int64(len(dPools))))
		stat.AVGScore /= int64(len(dPools))
	}

	stat.ActiveStake.SetLamports(ActiveStakeSum)
	stat.TotalSupply.SetLamports(SupplySum)
	stat.UnstakeLiquidity.SetLamports(UnstakeSum)

	s.cache.SetCurrentStatistic(stat, time.Second*30)

	return stat, nil
}

func (s *Imp) GetPoolStatistic(name string, aggregate string) ([]*smodels.Pool, error) {
	pool, err := s.dao.GetPool(name)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, fmt.Errorf("dao.GetPool(%s): %w", name, postgres.ErrorRecordNotFounded)
	}

	a, err := s.dao.GetPoolStatistic(pool.ID, postgres.SearchAggregate(aggregate))
	if err != nil {
		return nil, err
	}

	coin, err := s.dao.GetCoinByID(pool.CoinID)
	if err != nil {
		return nil, fmt.Errorf("dao.GetCoinByID: %w", err)
	}

	data := make([]*smodels.Pool, len(a))
	for i, v := range a {
		data[i] = (&smodels.Pool{}).Set(v, coin, pool, nil)
		data[i].ValidatorCount, err = s.dao.GetValidatorCount(&postgres.PoolValidatorDataCondition{
			PoolDataIDs: []uuid.UUID{
				v.ID,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (s *Imp) GetPoolCount() (int64, error) {
	i, err := s.dao.GetPoolCount(nil)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (s *Imp) GetNetworkAPY() (float64, error) {
	d, err := s.cache.GetAPY()
	if err != nil {
		return 0, err
	}

	f, _ := d.Float64()

	return f, nil
}

func (s *Imp) GetUSD() (float64, error) {
	d, err := s.cache.GetPrice()
	if err != nil {
		return 0, err
	}

	f, _ := d.Float64()

	return f, nil
}
