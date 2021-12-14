package dmodels

import (
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

type PoolValidatorData struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	PoolDataID  uuid.UUID `gorm:"type:uuid;not null;"`
	ValidatorID string    `gorm:"type:varchar(44);not null;"`
	ActiveStake uint64    `gorm:"type:int;not null;"`
	CreatedAt   time.Time `gorm:"index;not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

type Validator struct {
	ID              string          `gorm:"primaryKey;type:varchar(44);not null;"`
	Image           string          `gorm:"type:text"`
	Name            string          `gorm:"type:varchar(100);not null;"`
	Delinquent      bool            `gorm:"not null"`
	Network         string          `gorm:"type:varchar(50);not null;"`
	VotePK          string          `gorm:"index:idx_vote_pk,unique;type:varchar(44);not null;"`
	APY             decimal.Decimal `gorm:"type:decimal(8,4);not null;"`
	StakingAccounts uint64          `gorm:"type:int;not null;"`
	ActiveStake     uint64          `gorm:"type:int;not null;"`
	Fee             decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	Score           int64           `gorm:"type:int;not null;"`
	SkippedSlots    decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	DataCenter      string          `gorm:"not null"`
	CreatedAt       time.Time       `gorm:"index;not null"`
	UpdatedAt       time.Time       `gorm:"not null"`
}
