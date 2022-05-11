package services_test

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"gotest.tools/assert"
	"testing"
)

func TestGetLiquidityPools(t *testing.T) {
	data := map[string]struct {
		DAO  services.Imp
		Data struct {
			name   string
			limit  uint64
			offset uint64
		}
		Result []*smodels.LiquidityPool
		Err    error
	}{
		"first": {
			Data: struct {
				name   string
				limit  uint64
				offset uint64
			}{name: "LPName1", limit: 10, offset: 0},
			Result: []*smodels.LiquidityPool{
				{
					Name:  "LPName1",
					About: "123fdg",
					Image: "LPImg1",
					URL:   "LPUrl1",
				},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetLiquidityPoolsFunc: func(cond *postgres.Condition) ([]*dmodels.LiquidityPool, error) {
						if cond.Name != "LPName1" {
							return nil, fmt.Errorf("cond.Name != LPName1, name is %s", cond.Name)
						}
						if cond.Pagination.Limit != 10 {
							return nil, fmt.Errorf("limit != 10, limit = %d", cond.Pagination.Limit)
						}
						if cond.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Pagination.Offset)
						}
						return []*dmodels.LiquidityPool{&LPArr[0]}, nil
					},
					GetLiquidityPoolsCountFunc: func(cond *postgres.Condition) (int64, error) {
						if cond.Name != "LPName1" {
							return 0, fmt.Errorf("cond.Name != LPName1, name is %s", cond.Name)
						}
						return 1, nil
					},
				},
			},
		},
		"second": {
			Data: struct {
				name   string
				limit  uint64
				offset uint64
			}{name: "LPName1", limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetLiquidityPools: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetLiquidityPoolsFunc: func(cond *postgres.Condition) ([]*dmodels.LiquidityPool, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"third": {
			Data: struct {
				name   string
				limit  uint64
				offset uint64
			}{name: "LPName1", limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetLiquidityPoolsCount: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetLiquidityPoolsFunc: func(cond *postgres.Condition) ([]*dmodels.LiquidityPool, error) {
						if cond.Name != "LPName1" {
							return nil, fmt.Errorf("cond.Name != LPName1, name is %s", cond.Name)
						}
						if cond.Pagination.Limit != 10 {
							return nil, fmt.Errorf("limit != 10, limit = %d", cond.Pagination.Limit)
						}
						if cond.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Pagination.Offset)
						}
						return []*dmodels.LiquidityPool{&LPArr[0]}, nil
					},
					GetLiquidityPoolsCountFunc: func(cond *postgres.Condition) (int64, error) {
						return 0, fmt.Errorf("some error")
					},
				},
			},
		},
	}

	for s, s2 := range data {
		t.Run(s, func(t *testing.T) {
			lp, count, err := s2.DAO.GetLiquidityPools(s2.Data.name, s2.Data.limit, s2.Data.offset)
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			assert.Equal(t, uint64(len(s2.Result)), count)
			assert.Equal(t, uint64(len(lp)), count)
			for i, coin := range lp {
				t.Run(fmt.Sprintf("LiquidityPools[%d]", i), func(t *testing.T) {
					assert.DeepEqual(t, coin, s2.Result[i])
				})
			}
		})
	}
}
