package services_test

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/everstake/solana-pools/internal/services/smodels"
	uuid "github.com/satori/go.uuid"
	"gotest.tools/assert"
	"testing"
)

var GArr = []*dmodels.Governance{
	{
		ID:                 uuid.Must(uuid.FromString("721dd49b-0e19-4655-9052-42c8a57aef01")),
		Name:               "gov1",
		Symbol:             "g1",
		VoteURL:            "someurl",
		WebSiteURL:         "wbsturl",
		Image:              "img1",
		GeckoKey:           "key1",
		Blockchain:         "bc1",
		ContractAddress:    "addr1",
		MaximumTokenSupply: 1000000,
		CirculatingSupply:  1000000,
		USD:                85,
	},
	{
		ID:                 uuid.Must(uuid.FromString("1e85fd6d-3d32-4d86-b9a0-5ca2f7260af9")),
		Name:               "gov2",
		Symbol:             "g2",
		VoteURL:            "someurl2",
		WebSiteURL:         "wbsturl2",
		Image:              "img2",
		GeckoKey:           "key2",
		Blockchain:         "bc2",
		ContractAddress:    "addr2",
		MaximumTokenSupply: 2000000,
		CirculatingSupply:  2000000,
		USD:                95,
	},
}

func TestGetGovernance(t *testing.T) {
	data := map[string]struct {
		DAO  services.Imp
		Data struct {
			name   string
			sort   string
			desc   bool
			limit  uint64
			offset uint64
		}
		Result []*smodels.Governance
		Err    error
	}{
		"first": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "gov1", sort: "price", desc: true, limit: 10, offset: 0},
			Result: []*smodels.Governance{
				{
					Name:               "gov1",
					Symbol:             "g1",
					VoteURL:            "someurl",
					WebSiteURL:         "wbsturl",
					Image:              "img1",
					GeckoKey:           "key1",
					Blockchain:         "bc1",
					ContractAddress:    "addr1",
					MaximumTokenSupply: 1000000,
					CirculatingSupply:  1000000,
					USD:                85,
				},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetGovernanceFunc: func(cond *postgres.GovernanceCondition) ([]*dmodels.Governance, error) {
						if cond.Condition.Name != "gov1" {
							return nil, fmt.Errorf("cond.Condition.Name != gov1, name is %s", cond.Condition.Name)
						}
						if cond.Condition.Pagination.Limit != 10 {
							return nil, fmt.Errorf("limit != 10, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Condition.Pagination.Offset)
						}
						if cond.Sort.Sort != postgres.GovernancePrice {
							return nil, fmt.Errorf("cond.Sort.Sort != %d, sort = %d", postgres.GovernancePrice, cond.Sort.Sort)
						}
						if cond.Sort.Desc != true {
							return nil, fmt.Errorf("cond.Sort.Desc != true, desc = %v", false)
						}
						return GArr[:1], nil
					},
					GetGovernanceCountFunc: func(cond *postgres.GovernanceCondition) (int64, error) {
						if cond.Condition.Name != "gov1" {
							return 0, fmt.Errorf("cond.Condition.Name != gov1, name is %s", cond.Condition.Name)
						}
						return 1, nil
					},
				},
			},
		},
		"second": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "gov3", sort: "price", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetGovernanceFunc: func(cond *postgres.GovernanceCondition) ([]*dmodels.Governance, error) {
						if cond.Condition.Name != "gov3" {
							return nil, fmt.Errorf("cond.Condition.Name != gov1, name is %s", cond.Condition.Name)
						}
						if cond.Condition.Pagination.Limit != 10 {
							return nil, fmt.Errorf("limit != 10, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Condition.Pagination.Offset)
						}
						if cond.Sort.Sort != postgres.GovernancePrice {
							return nil, fmt.Errorf("cond.Sort.Sort != %d, sort = %d", postgres.GovernancePrice, cond.Sort.Sort)
						}
						if cond.Sort.Desc != true {
							return nil, fmt.Errorf("cond.Sort.Desc != true, desc = %v", false)
						}
						return nil, nil
					},
					GetGovernanceCountFunc: func(cond *postgres.GovernanceCondition) (int64, error) {
						if cond.Condition.Name != "gov3" {
							return 0, fmt.Errorf("cond.Condition.Name != gov1, name is %s", cond.Condition.Name)
						}
						return 0, nil
					},
				},
			},
		},
		"third": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "gov2", sort: "price", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetCoinsCount: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetGovernanceFunc: func(cond *postgres.GovernanceCondition) ([]*dmodels.Governance, error) {
						if cond.Condition.Name != "gov2" {
							return nil, fmt.Errorf("cond.Condition.Name != gov1, name is %s", cond.Condition.Name)
						}
						if cond.Condition.Pagination.Limit != 10 {
							return nil, fmt.Errorf("limit != 10, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Condition.Pagination.Offset)
						}
						if cond.Sort.Sort != postgres.GovernancePrice {
							return nil, fmt.Errorf("cond.Sort.Sort != %d, sort = %d", postgres.GovernancePrice, cond.Sort.Sort)
						}
						if cond.Sort.Desc != true {
							return nil, fmt.Errorf("cond.Sort.Desc != true, desc = %v", false)
						}
						return GArr[1:], nil
					},
					GetGovernanceCountFunc: func(cond *postgres.GovernanceCondition) (int64, error) {
						if cond.Condition.Name == "gov2" {
							return 0, fmt.Errorf("some error")
						}
						return 1, nil
					},
				},
			},
		},
	}

	for s, s2 := range data {
		t.Run(s, func(t *testing.T) {
			gov, count, err := s2.DAO.GetGovernance(s2.Data.name, s2.Data.sort, s2.Data.desc, s2.Data.limit, s2.Data.offset)
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			assert.Equal(t, uint64(len(s2.Result)), count)
			assert.Equal(t, uint64(len(gov)), count)
			for i, coin := range gov {
				t.Run(fmt.Sprintf("governances[%d]", i), func(t *testing.T) {
					assert.DeepEqual(t, coin, s2.Result[i])
				})
			}
		})
	}
}
