package dmodels

import (
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Validator struct {
	ID            uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	PoolDataID    uuid.UUID       `gorm:"type:uuid;not null;"`
	NodePK        string          `gorm:"primaryKey, type:varchar(44);not null;"`
	APY           decimal.Decimal `gorm:"type:decimal(8,4);not null;"`
	VotePK        string          `gorm:"index, type:varchar(44);not null;"`
	StakeAccounts uint64          `gorm:"type:int;not null;"`
	ActiveStake   decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	Fee           decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	Score         int64           `gorm:"type:int;not null;"`
	SkippedSlots  decimal.Decimal `gorm:"type:decimal(5,2);not null;"`
	DataCenter    string          `gorm:"not null"`
	CreatedAt     time.Time       `gorm:"index;not null"`
	UpdatedAt     time.Time       `gorm:"not null"`
}
