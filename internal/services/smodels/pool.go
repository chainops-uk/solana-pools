package smodels

import (
	"github.com/shopspring/decimal"
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
	}
	PoolDetails struct {
		Pool
		Validators []*Validator
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
