package services

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

func (s *Imp) GetPool(name string) (pool smodels.PoolDetails, err error) {
	dPool, err := s.dao.GetPool(name)
	if err != nil {
		return pool, fmt.Errorf("dao.GetPool: %s", err.Error())
	}
	dLastPoolData, err := s.dao.GetLastPoolData(dPool.ID, nil)
	if err != nil {
		return pool, fmt.Errorf("dao.GetPoolData: %s", err.Error())
	}
	dValidators, err := s.dao.GetValidators(dPool.ID)
	if err != nil {
		return pool, fmt.Errorf("dao.GetValidators: %s", err.Error())
	}
	validators := make([]*smodels.Validator, len(dValidators))
	for i, v := range dValidators {
		validators[i] = &smodels.Validator{
			NodePK:       v.NodePK,
			APY:          v.APY,
			VotePK:       v.VotePK,
			ActiveStake:  v.ActiveStake,
			Fee:          v.Fee,
			Score:        v.Score,
			SkippedSlots: v.SkippedSlots,
			DataCenter:   v.DataCenter,
		}
	}
	return smodels.PoolDetails{
		Pool: smodels.Pool{
			Address:          dPool.Address,
			Name:             dPool.Name,
			ActiveStake:      dLastPoolData.ActiveStake,
			TokensSupply:     dLastPoolData.TotalTokensSupply,
			APY:              dLastPoolData.APY,
			AVGSkippedSlots:  dLastPoolData.AVGSkippedSlots,
			AVGScore:         dLastPoolData.AVGScore,
			Delinquent:       dLastPoolData.Delinquent,
			UnstakeLiquidity: dLastPoolData.UnstakeLiquidity,
			DepossitFee:      dLastPoolData.DepossitFee,
			WithdrawalFee:    dLastPoolData.WithdrawalFee,
			RewardsFee:       dLastPoolData.RewardsFee,
		},
		Validators: validators,
	}, nil
}

func (s *Imp) GetPools(name string, limit uint64, offset uint64) ([]*smodels.PoolDetails, error) {
	dPools, err := s.dao.GetPools(&postgres.Condition{
		Name: name,
		Pagination: postgres.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("dao.GetPool: %w", err)
	}
	if len(dPools) == 0 {
		return nil, nil
	}
	pools := make([]*smodels.PoolDetails, len(dPools))
	for i, v1 := range dPools {
		pools[i] = &smodels.PoolDetails{
			Pool: smodels.Pool{
				Address: v1.Address,
				Name:    v1.Name,
			},
		}

		dLastPoolData, err := s.dao.GetLastPoolData(v1.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("dao.GetLastPoolData: %w", err)
		}

		pools[i].Set(dLastPoolData, &v1)

		pools[i].ActiveStake = dLastPoolData.ActiveStake
		pools[i].TokensSupply = dLastPoolData.TotalTokensSupply
		pools[i].APY = dLastPoolData.APY
		pools[i].AVGSkippedSlots = dLastPoolData.AVGSkippedSlots
		pools[i].AVGScore = dLastPoolData.AVGScore
		pools[i].Delinquent = dLastPoolData.Delinquent
		pools[i].UnstakeLiquidity = dLastPoolData.UnstakeLiquidity
		pools[i].DepossitFee = dLastPoolData.DepossitFee
		pools[i].WithdrawalFee = dLastPoolData.WithdrawalFee
		pools[i].RewardsFee = dLastPoolData.RewardsFee

		dValidators, err := s.dao.GetValidators(dLastPoolData.ID)
		if err != nil {
			return nil, fmt.Errorf("dao.GetValidators: %w", err)
		}

		validators := make([]*smodels.Validator, len(dValidators))
		for i, v2 := range dValidators {
			validators[i] = &smodels.Validator{
				NodePK:       v2.NodePK,
				APY:          v2.APY,
				VotePK:       v2.VotePK,
				ActiveStake:  v2.ActiveStake,
				Fee:          v2.Fee,
				Score:        v2.Score,
				SkippedSlots: v2.SkippedSlots,
				DataCenter:   v2.DataCenter,
			}
		}

		pools[i].Validators = validators
	}

	return pools, nil
}

func (s *Imp) GetPoolsCurrentStatistic() (*smodels.Statistic, error) {
	dPools, err := s.dao.GetPools(nil)
	if err != nil {
		return nil, fmt.Errorf("dao.GetPool: %w", err)
	}
	if len(dPools) == 0 {
		return nil, nil
	}
	stat := &smodels.Statistic{}

	once := sync.Once{}
	for _, v := range dPools {
		dLastPoolData, err := s.dao.GetLastPoolData(v.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("dao.GetLastPoolData: %w", err)
		}
		if dLastPoolData == nil {
			continue
		}

		once.Do(func() {
			stat.MINScore = dLastPoolData.AVGScore
			stat.MAXScore = dLastPoolData.AVGScore
		})

		if dLastPoolData.AVGScore > stat.MAXScore {
			stat.MAXScore = dLastPoolData.AVGScore
		}
		if dLastPoolData.AVGScore < stat.MINScore {
			stat.MINScore = dLastPoolData.AVGScore
		}

		stat.ActiveStake = stat.ActiveStake.Add(dLastPoolData.ActiveStake)
		stat.AVGSkippedSlots = stat.AVGSkippedSlots.Add(dLastPoolData.AVGSkippedSlots)
		stat.AVGScore += dLastPoolData.AVGScore
		stat.Delinquent = stat.Delinquent.Add(dLastPoolData.Delinquent)
		stat.UnstakeLiquidity = stat.UnstakeLiquidity.Add(dLastPoolData.UnstakeLiquidity)
	}

	stat.AVGSkippedSlots = stat.AVGSkippedSlots.Div(decimal.NewFromInt(int64(len(dPools))))
	stat.AVGScore /= int64(len(dPools))

	return stat, nil
}

func (s *Imp) GetPoolsStatistic(name string, aggregate string, from time.Time, to time.Time) ([]*smodels.Pool, error) {
	pool, err := s.dao.GetPool(name)
	if err != nil {
		return nil, err
	}

	a, err := s.dao.GetPoolStatistic(pool.ID, postgres.SearchAggregate(aggregate), from, to)
	if err != nil {
		return nil, err
	}

	data := make([]*smodels.Pool, len(a))
	for i, v := range a {
		data[i] = (&smodels.Pool{}).Set(v, &pool)
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
