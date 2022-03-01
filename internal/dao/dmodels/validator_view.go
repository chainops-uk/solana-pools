package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

type ValidatorView struct {
	ID              string          `gorm:"column:id"`
	Image           string          `gorm:"column:image"`
	Name            string          `gorm:"column:name"`
	Delinquent      bool            `gorm:"column:delinquent"`
	VotePK          string          `gorm:"column:vote_pk"`
	APY             decimal.Decimal `gorm:"column:apy"`
	StakingAccounts uint64          `gorm:"column:staking_accounts"`
	ActiveStake     uint64          `gorm:"column:active_stake"`
	Fee             decimal.Decimal `gorm:"column:fee"`
	Score           int64           `gorm:"column:score"`
	SkippedSlots    decimal.Decimal `gorm:"column:skipped_slots"`
	DataCenter      string          `gorm:"column:data_center"`
	CreatedAt       time.Time       `gorm:"column:created_at"`
	UpdatedAt       time.Time       `gorm:"column:updated_at"`
}
