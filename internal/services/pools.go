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
	dValidators, err := s.dao.GetValidators(dPool.ID)
	if err != nil {
		return pool, fmt.Errorf("dao.GetValidators: %s", err.Error())
	}
	validators := make([]smodels.Validator, len(dValidators))
	for i, v := range dValidators {
		validators[i] = smodels.Validator{
			NodePK:       v.NodePK,
			APR:          v.APR,
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
			Name:             dPool.Address,
			ActiveStake:      dPool.ActiveStake,
			TokensSupply:     dPool.TokensSupply,
			APR:              dPool.APR,
			Nodes:            dPool.Nodes,
			AVGSkippedSlots:  dPool.AVGSkippedSlots,
			AVGScore:         dPool.AVGScore,
			Delinquent:       dPool.Delinquent,
			UnstakeLiquidity: dPool.UnstakeLiquidity,
			DepossitFee:      dPool.DepossitFee,
			WithdrawalFee:    dPool.WithdrawalFee,
			RewardsFee:       dPool.RewardsFee,
		},
		Validators: validators,
	}, nil
}
