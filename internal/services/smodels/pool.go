package smodels

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/pkg/models/sol"
	"github.com/shopspring/decimal"
	"time"
)

type (
	Pool struct {
		Address          string
		Name             string
		Image            string
		Currency         string
		ActiveStake      sol.SOL
		TokensSupply     sol.SOL
		TotalLamports    sol.SOL
		APY              decimal.Decimal
		AVGSkippedSlots  decimal.Decimal
		AVGScore         int64
		StakingAccounts  uint64
		Delinquent       decimal.Decimal
		UnstakeLiquidity sol.SOL
		DepossitFee      decimal.Decimal
		WithdrawalFee    decimal.Decimal
		RewardsFee       decimal.Decimal
		ValidatorCount   int64
		CreatedAt        time.Time
	}
	PoolDetails struct {
		Pool
		CreatedAt time.Time
	}
	Statistic struct {
		Pools            uint64
		ActiveStake      sol.SOL
		TotalSupply      sol.SOL
		AVGSkippedSlots  decimal.Decimal
		AVGPoolsApy      decimal.Decimal
		MAXScore         int64
		AVGScore         int64
		MINScore         int64
		Delinquent       decimal.Decimal
		UnstakeLiquidity sol.SOL
	}
)

func (p *Pool) Set(data *dmodels.PoolData, coin *dmodels.Coin, pool *dmodels.Pool, validator []*dmodels.Validator) *Pool {
	if pool != nil {
		p.Name = pool.Name
		p.Address = pool.Address
		p.Image = pool.Image
	}
	if coin != nil {
		p.Currency = coin.Name
	}
	if data != nil {
		p.ActiveStake.SetLamports(data.ActiveStake)
		p.TokensSupply.SetLamports(data.TotalTokensSupply)
		p.TotalLamports.SetLamports(data.TotalLamports)
		p.APY = data.APY

		p.UnstakeLiquidity.SetLamports(data.UnstakeLiquidity)
		p.DepossitFee = data.DepossitFee
		p.WithdrawalFee = data.WithdrawalFee
		p.RewardsFee = data.RewardsFee
		p.CreatedAt = data.CreatedAt
	}
	if validator != nil {
		p.ValidatorCount = int64(len(validator))
		if len(validator) > 0 {
			for _, v := range validator {
				p.AVGScore += v.Score
				p.AVGSkippedSlots = p.AVGSkippedSlots.Add(v.SkippedSlots)
				if v.Delinquent {
					p.Delinquent = p.Delinquent.Add(decimal.NewFromInt(1))
				}
				p.StakingAccounts += v.StakingAccounts
			}
			p.AVGSkippedSlots = p.AVGSkippedSlots.Div(decimal.NewFromInt(int64(len(validator))))
			p.AVGScore /= int64(len(validator))
			p.Delinquent = p.Delinquent.Div(decimal.NewFromInt(int64(len(validator))))

		}
	}
	return p
}
