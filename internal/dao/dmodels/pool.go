package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

type Pool struct {
	ID               uint64          `gorm:"primaryKey, type:not null;"`
	Address          string          `gorm:"index, not null;"`
	Network          string          `gorm:"type:varchar(50);not null;"`
	Active           bool            `gorm:"not null"`
	Name             string          `gorm:"type:varchar(100);not null;"`
	ActiveStake      decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	TokensSupply     decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	APR              decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	Nodes            uint64          `gorm:"type:int;not null;"`
	AVGSkippedSlots  decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	AVGScore         int64           `gorm:"type:int;not null;"`
	Delinquent       decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	UnstakeLiquidity decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	DepossitFee      decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	WithdrawalFee    decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	RewardsFee       decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	UpdatedAt        time.Time       `gorm:"not null"`
	CreatedAt        time.Time       `gorm:"not null"`
}
