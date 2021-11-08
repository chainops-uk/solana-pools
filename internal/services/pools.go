package services

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/services/smodels"
)

func (s Imp) GetPool(name string) (pool smodels.PoolDetails, err error) {
	dPool, err := s.dao.GetPool(name)
	if err != nil {
		return pool, fmt.Errorf("dao.GetPool: %s", err.Error())
	}
	dLastPoolData, err := s.dao.GetLastPoolData(dPool.ID)
	if err != nil {
		return pool, fmt.Errorf("dao.GetPoolData: %s", err.Error())
	}
	dValidators, err := s.dao.GetValidators(dPool.ID)
	if err != nil {
		return pool, fmt.Errorf("dao.GetValidators: %s", err.Error())
	}
	validators := make([]smodels.Validator, len(dValidators))
	for i, v := range dValidators {
		validators[i] = smodels.Validator{
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
			TokensSupply:     dLastPoolData.TokensSupply,
			APR:              dLastPoolData.APR,
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
