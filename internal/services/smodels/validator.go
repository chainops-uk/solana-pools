package smodels

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/pkg/models/sol"
	"github.com/shopspring/decimal"
)

type Validator struct {
	Name             string
	Image            string
	StakingAccounts  uint64
	NodePK           string
	APY              decimal.Decimal
	VotePK           string
	PoolActiveStake  sol.SOL
	TotalActiveStake sol.SOL
	Fee              decimal.Decimal
	Score            int64
	SkippedSlots     decimal.Decimal
	DataCenter       string
}

func (v *Validator) Set(activeStake uint64, vv *dmodels.Validator) *Validator {
	if vv == nil {
		v.PoolActiveStake.SetLamports(activeStake)
		return v
	}
	v.Name = vv.Name
	v.Image = vv.Image
	v.StakingAccounts = vv.StakingAccounts
	v.NodePK = vv.ID
	v.APY = vv.APY
	v.VotePK = vv.VotePK
	v.PoolActiveStake.SetLamports(activeStake)
	v.TotalActiveStake.SetLamports(vv.ActiveStake)
	v.Fee = vv.Fee
	v.Score = vv.Score
	v.SkippedSlots = vv.SkippedSlots
	v.DataCenter = vv.DataCenter
	return v
}
