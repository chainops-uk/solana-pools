package smodels

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/shopspring/decimal"
	"time"
)

type (
	Pool struct {
		Address          string
		Name             string
		ActiveStake      decimal.Decimal
		TokensSupply     decimal.Decimal
		APY              decimal.Decimal
		AVGSkippedSlots  decimal.Decimal
		AVGScore         int64
		Delinquent       decimal.Decimal
		UnstakeLiquidity decimal.Decimal
		DepossitFee      decimal.Decimal
		WithdrawalFee    decimal.Decimal
		RewardsFee       decimal.Decimal
		CreatedAt        time.Time
	}
	PoolDetails struct {
		Pool
		Validators []*Validator
		CreatedAt  time.Time
	}
	Statistic struct {
		ActiveStake      decimal.Decimal
		AVGSkippedSlots  decimal.Decimal
		MAXScore         int64
		AVGScore         int64
		MINScore         int64
		Delinquent       decimal.Decimal
		UnstakeLiquidity decimal.Decimal
	}
)

func (p *Pool) Set(data *dmodels.PoolData, pool *dmodels.Pool) *Pool {
	if pool != nil {
		p.Name = pool.Name
		p.Address = pool.Address
	}
	if data != nil {
		p.ActiveStake = data.ActiveStake
		p.TokensSupply = data.TotalTokensSupply
		p.APY = data.APY
		p.AVGSkippedSlots = data.AVGSkippedSlots
		p.AVGScore = data.AVGScore
		p.Delinquent = data.Delinquent
		p.UnstakeLiquidity = data.UnstakeLiquidity
		p.DepossitFee = data.DepossitFee
		p.WithdrawalFee = data.WithdrawalFee
		p.RewardsFee = data.RewardsFee
		p.CreatedAt = data.CreatedAt
	}
	return p
}
