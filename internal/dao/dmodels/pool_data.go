package dmodels

import (
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

type PoolData struct {
	ID                uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	PoolID            uuid.UUID       `gorm:"type:uuid;not null;"`
	Epoch             uint64          `gorm:"type:int8;not null;"`
	ActiveStake       uint64          `gorm:"type:int;not null;"`
	TotalTokensSupply uint64          `gorm:"type:int;not null;"`
	TotalLamports     uint64          `gorm:"type:int;not null;"`
	APY               decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	UnstakeLiquidity  uint64          `gorm:"type:int;not null;"`
	DepossitFee       decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	WithdrawalFee     decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	RewardsFee        decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	UpdatedAt         time.Time       `gorm:"not null"`
	CreatedAt         time.Time       `gorm:"index;not null"`
	Pool              Pool            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:Restrict;"`
}
