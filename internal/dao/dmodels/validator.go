package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

type Validator struct {
	NodePK       string          `gorm:"primaryKey, type:varchar(44);not null;"`
	PoolID       uint64          `gorm:"type:int;not null;"`
	APR          decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	VotePK       string          `gorm:"index, type:varchar(44);not null;"`
	ActiveStake  decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	Fee          decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	Score        int64           `gorm:"type:int;not null;"`
	SkippedSlots decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	DataCenter   string          `gorm:"not null"`
	UpdatedAt    time.Time       `gorm:"not null"`
	CreatedAt    time.Time       `gorm:"not null"`
}
