package smodels

import (
	"github.com/shopspring/decimal"
)

type Validator struct {
	NodePK       string
	APY          decimal.Decimal
	VotePK       string
	ActiveStake  decimal.Decimal
	Fee          decimal.Decimal
	Score        int64
	SkippedSlots decimal.Decimal
	DataCenter   string
}
