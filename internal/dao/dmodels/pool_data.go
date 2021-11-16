package dmodels

import (
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

type PoolData struct {
	ID               uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	PoolID           uuid.UUID       `gorm:"type:uuid;not null;"`
	ActiveStake      decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	TokensSupply     decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	APY              decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	AVGSkippedSlots  decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	AVGScore         int64           `gorm:"type:int;not null;"`
	Delinquent       decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	UnstakeLiquidity decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	DepossitFee      decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	WithdrawalFee    decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	RewardsFee       decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	UpdatedAt        time.Time       `gorm:"not null"`
	CreatedAt        time.Time       `gorm:"index;not null"`
}
