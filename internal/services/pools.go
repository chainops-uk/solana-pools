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

func (s *Imp) GetPool(name string, epoch uint64) (*smodels.PoolDetails, error) {
	pd, err := s.Cache.GetPool(name)
	if err != nil && !errors.Is(err, cache.KeyWasNotFound) {
		return nil, err
	}
	if pd != nil {
		return pd, nil
	}

	dPool, err := s.DAO.GetPool(name)
	if err != nil {
		return nil, fmt.Errorf("DAO.GetPool: %s", err.Error())
	}
	if dPool == nil {
		return nil, fmt.Errorf("DAO.GetPool(%s): %w", name, postgres.ErrorRecordNotFounded)
	}

	var dLastPoolData *dmodels.PoolData
	if epoch == 10 {
		dLastPoolData, err = s.DAO.GetLastPoolDataWithApyForTenEpoch(dPool.ID)
		if err != nil {
			return nil, fmt.Errorf("DAO.GetLastPoolDataWithApyForTenEpoch: %w", err)
		}
	} else {
		dLastPoolData, err = s.DAO.GetLastPoolData(dPool.ID)
		if err != nil {
			return nil, fmt.Errorf("DAO.GetLastPoolData: %w", err)
		}
	}

	dValidators, err := s.DAO.GetPoolValidatorData(&postgres.PoolValidatorDataCondition{PoolDataIDs: []uuid.UUID{dLastPoolData.ID}}, epoch)
	if err != nil {
		return nil, fmt.Errorf("DAO.GetPoolValidatorData: %s", err.Error())
	}
	validatorsS := make([]*smodels.PoolValidatorData, len(dValidators))
	validatorsD := make([]*dmodels.ValidatorView, len(dValidators))
	for i, v := range dValidators {
		validatorsD[i], err = s.DAO.GetValidator(v.ValidatorID, epoch)
		if err != nil {
			return nil, fmt.Errorf("DAO.GetValidator(%s): %w", v.ValidatorID, err)
		}
		validatorsS[i] = (&smodels.PoolValidatorData{}).Set(v.ActiveStake, validatorsD[i])
	}

	coin, err := s.DAO.GetCoinByID(dPool.CoinID)
	if err != nil {
		return nil, fmt.Errorf("DAO.GetCoinByID: %w", err)
	}
	Pool := (&smodels.Pool{}).Set(dLastPoolData, coin, dPool, validatorsD)

	pd = &smodels.PoolDetails{
		Pool: *Pool,
	}

	s.Cache.SetPool(pd, time.Second*30)

	return pd, nil
}

func (s *Imp) GetPools(name string, sort string, desc bool, epoch uint64, limit uint64, offset uint64) ([]*smodels.PoolDetails, uint64, error) {
	dPools, err := s.DAO.GetPools(&postgres.PoolCondition{
		Condition: &postgres.Condition{
			Network: postgres.MainNet,
			Name:    name,
			Pagination: postgres.Pagination{
				Limit:  limit,
				Offset: offset,
			},
		},
		Sort: &postgres.PoolDataSort{
			Epoch:    epoch,
			PoolSort: postgres.SearchPoolSort(sort),
			Desc:     desc,
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetPool: %w", err)
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

		var dLastPoolData *dmodels.PoolData
		if epoch == 10 {
			dLastPoolData, err = s.DAO.GetLastPoolDataWithApyForTenEpoch(v1.ID)
			if err != nil {
				return nil, 0, fmt.Errorf("DAO.GetLastPoolDataWithApyForTenEpoch: %w", err)
			}
		} else {
			dLastPoolData, err = s.DAO.GetLastPoolData(v1.ID)
			if err != nil {
				return nil, 0, fmt.Errorf("DAO.GetLastPoolData: %w", err)
			}
		}

		dValidators, err := s.DAO.GetPoolValidatorData(&postgres.PoolValidatorDataCondition{PoolDataIDs: []uuid.UUID{dLastPoolData.ID}}, epoch)
		if err != nil {
			return nil, 0, fmt.Errorf("DAO.GetPoolValidatorData: %w", err)
		}

		validatorsD := make([]*dmodels.ValidatorView, len(dValidators))
		for i, v2 := range dValidators {

			validatorsD[i], err = s.DAO.GetValidator(v2.ValidatorID, epoch)

			if err != nil {
				return nil, 0, fmt.Errorf("DAO.GetValidator: %w", err)
			}
		}

		coin, err := s.DAO.GetCoinByID(v1.CoinID)
		if err != nil {
			return nil, 0, fmt.Errorf("DAO.GetCoinByID: %w", err)
		}

		pools[i].Set(dLastPoolData, coin, v1, validatorsD)
	}

	count, err := s.DAO.GetPoolCount(&postgres.Condition{
		Network: postgres.MainNet,
		Name:    name,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetPoolCount: %w", err)
	}

	return pools, uint64(count), nil
}

func (s *Imp) GetPoolsCurrentStatistic(epoch uint64) (*smodels.Statistic, error) {
	stat, err := s.Cache.GetCurrentStatistic()

	if err != nil && !errors.Is(err, cache.KeyWasNotFound) {
		return nil, err
	}
	if stat != nil {
		return stat, nil
	}

	dPools, err := s.DAO.GetPools(&postgres.PoolCondition{Condition: &postgres.Condition{Network: postgres.MainNet}})
	if err != nil {
		return nil, fmt.Errorf("DAO.GetPool: %w", err)
	}
	if len(dPools) == 0 {
		return nil, nil
	}
	stat = &smodels.Statistic{
		Pools:       uint64(len(dPools)),
		MAXPoolsApy: decimal.NewFromInt(0),
	}

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

		var dLastPoolData *dmodels.PoolData
		if epoch == 10 {
			dLastPoolData, err = s.DAO.GetLastPoolDataWithApyForTenEpoch(v1.ID)
			if err != nil {
				return nil, fmt.Errorf("DAO.GetLastPoolDataWithApyForTenEpoch: %w", err)
			}
		} else {
			dLastPoolData, err = s.DAO.GetLastPoolData(v1.ID)
			if err != nil {
				return nil, fmt.Errorf("DAO.GetLastPoolData: %w", err)
			}
		}

		dValidators, err := s.DAO.GetPoolValidatorData(&postgres.PoolValidatorDataCondition{PoolDataIDs: []uuid.UUID{dLastPoolData.ID}}, epoch)
		if err != nil {
			return nil, fmt.Errorf("DAO.GetValidators: %w", err)
		}

		validatorsD := make([]*dmodels.ValidatorView, len(dValidators))
		for i, v2 := range dValidators {

			validatorsD[i], err = s.DAO.GetValidator(v2.ValidatorID, epoch)
			if err != nil {
				return nil, fmt.Errorf("DAO.GetValidator: %w", err)
			}
		}

		coin, err := s.DAO.GetCoinByID(v1.CoinID)
		if err != nil {
			return nil, fmt.Errorf("DAO.GetCoinByID: %w", err)
		}

		pools[i].Set(dLastPoolData, coin, v1, validatorsD)

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
		stat.Delinquent += pools[i].Delinquent
		if pools[i].APY.GreaterThan(stat.MAXPoolsApy) {
			stat.MAXPoolsApy = pools[i].APY
		}
	}

	if len(dPools) > 0 {
		stat.AVGSkippedSlots = stat.AVGSkippedSlots.Div(decimal.NewFromInt(int64(len(dPools))))
		stat.AVGScore /= int64(len(dPools))
	}

	stat.ActiveStake.SetLamports(ActiveStakeSum)
	stat.TotalSupply.SetLamports(SupplySum)
	stat.UnstakeLiquidity.SetLamports(UnstakeSum)

	s.Cache.SetCurrentStatistic(stat, time.Second*30)

	return stat, nil
}

func (s *Imp) GetPoolStatistic(name string, aggregate string) ([]*smodels.Pool, error) {
	pool, err := s.DAO.GetPool(name)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, fmt.Errorf("DAO.GetPool(%s): %w", name, postgres.ErrorRecordNotFounded)
	}

	a, err := s.DAO.GetPoolStatistic(pool.ID, postgres.SearchAggregate(aggregate))
	if err != nil {
		return nil, err
	}

	coin, err := s.DAO.GetCoinByID(pool.CoinID)
	if err != nil {
		return nil, fmt.Errorf("DAO.GetCoinByID: %w", err)
	}

	data := make([]*smodels.Pool, len(a))
	for i, v := range a {
		data[i] = (&smodels.Pool{}).Set(v, coin, pool, nil)
		data[i].ValidatorCount, err = s.DAO.GetValidatorDataCount(&postgres.PoolValidatorDataCondition{
			PoolDataIDs: []uuid.UUID{
				v.ID,
			},
		}, 1)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (s *Imp) GetNetworkAPY() (float64, error) {
	d, err := s.Cache.GetAPY()
	if err != nil {
		return 0, err
	}

	f, _ := d.Float64()

	return f, nil
}

func (s *Imp) GetUSD() (float64, error) {
	d, err := s.Cache.GetPrice()
	if err != nil {
		return 0, err
	}

	f, _ := d.Float64()

	return f, nil
}
