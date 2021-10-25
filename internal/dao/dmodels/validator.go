package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

type Validator struct {
	PoolID       uint64
	APR          uint64
	VotePK       string
	NodePK       string
	ActiveStake  decimal.Decimal
	Fee          decimal.Decimal
	Score        int64
	SkippedSlots decimal.Decimal
	DataCenter   string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}
