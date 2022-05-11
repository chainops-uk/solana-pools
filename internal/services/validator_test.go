package services_test

import (
	"github.com/everstake/solana-pools/internal/dao"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/everstake/solana-pools/pkg/models/sol"
	"github.com/shopspring/decimal"
	"testing"
)

func TestGetPoolValidators(t *testing.T) {
	data := map[string]struct {
		DAO  services.Imp
		Data struct {
			name          string
			validatorName string
			sort          string
			desc          bool
			limit         uint64
			offset        uint64
		}
		Result []*smodels.PoolValidatorData
		Err    error
	}{
		"first": {
			Data: struct {
				name          string
				validatorName string
				sort          string
				desc          bool
				limit         uint64
				offset        uint64
			}{name: "pool1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: []*smodels.PoolValidatorData{
				{
					Name:             "",
					Image:            "",
					StakingAccounts:  0,
					NodePK:           "",
					APY:              decimal.Decimal{},
					VotePK:           "",
					PoolActiveStake:  sol.SOL{},
					TotalActiveStake: sol.SOL{},
					Fee:              decimal.Decimal{},
					Score:            0,
					SkippedSlots:     decimal.Decimal{},
					DataCenter:       "",
				},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{},
			},
		},
	}
}
