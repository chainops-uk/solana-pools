package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

type Pool struct {
	ID               uint64
	Address          string
	Network          string
	Active           bool
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
	UpdatedAt        time.Time
	CreatedAt        time.Time
}
