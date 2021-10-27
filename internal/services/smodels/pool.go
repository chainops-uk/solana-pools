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
		APR              decimal.Decimal
		Nodes            uint64
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
		Validators []Validator
	}
)
