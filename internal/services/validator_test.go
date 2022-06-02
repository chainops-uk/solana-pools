package services_test

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/everstake/solana-pools/pkg/models/sol"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"gotest.tools/assert"
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
			}{name: "pool1", validatorName: "val1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: []*smodels.PoolValidatorData{
				{
					Name:             "val1",
					Image:            "img1",
					StakingAccounts:  500,
					NodePK:           "id1",
					APY:              decimal.Decimal{},
					VotePK:           "pk1",
					PoolActiveStake:  sol.SOL{decimal.NewFromFloat(0.000854684)},
					TotalActiveStake: sol.SOL{decimal.NewFromFloat(0.0000001)},
					Fee:              decimal.Decimal{},
					Score:            5698,
					SkippedSlots:     decimal.Decimal{},
					DataCenter:       "dc",
				},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != dPool.ID, poolID is %s", poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition, epoch uint64) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						if condition.Sort.ValidatorDataSort != 1 {
							return nil, fmt.Errorf("condition.Sort.ValidatorDataSort != 1, condition.Sort.ValidatorDataSort is %d", condition.Sort.ValidatorDataSort)
						}
						if condition.Sort.Desc != true {
							return nil, fmt.Errorf("condition.Sort.Desc != true, condition.Sort.Desc is %v", false)
						}
						if condition.Condition.Name != "val1" {
							return nil, fmt.Errorf("condition.Condition.Name != val1, condition.Condition.Name is %s", condition.Condition.Name)
						}
						if condition.Condition.Pagination.Limit != 10 {
							return nil, fmt.Errorf("condition.Condition.Pagination.Limit != 10, but %d", condition.Condition.Pagination.Limit)
						}
						if condition.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("condition.Condition.Pagination.Offset != 10, but %d", condition.Condition.Pagination.Offset)
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string, epoch uint64) (*dmodels.ValidatorView, error) {
						if validatorID != poolVD[0].ValidatorID {
							return nil, fmt.Errorf("validatorID != %s, validatorID is %s", poolVD[0].ValidatorID, validatorID)
						}
						return &dValView, nil
					},
					GetValidatorDataCountFunc: func(condition *postgres.PoolValidatorDataCondition, epoch uint64) (int64, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return 0, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						if condition.Condition.Name != "val1" {
							return 0, fmt.Errorf("condition.Condition.Name != val1, condition.Condition.Name is %s", condition.Condition.Name)
						}
						return 1, nil
					},
				},
			},
		},
		"second": {
			Data: struct {
				name          string
				validatorName string
				sort          string
				desc          bool
				limit         uint64
				offset        uint64
			}{name: "pool1", validatorName: "val1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPool(%s): %w", "pool1", postgres.ErrorRecordNotFounded),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						return nil, nil
					},
				},
			},
		},
		"third": {
			Data: struct {
				name          string
				validatorName string
				sort          string
				desc          bool
				limit         uint64
				offset        uint64
			}{name: "pool1", validatorName: "val1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPool: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"forth": {
			Data: struct {
				name          string
				validatorName string
				sort          string
				desc          bool
				limit         uint64
				offset        uint64
			}{name: "pool1", validatorName: "val1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetLastPoolData: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"fifth": {
			Data: struct {
				name          string
				validatorName string
				sort          string
				desc          bool
				limit         uint64
				offset        uint64
			}{name: "pool1", validatorName: "val1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPoolValidatorData: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != dPool.ID, poolID is %s", poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition, epoch uint64) ([]*dmodels.PoolValidatorData, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"sixth": {
			Data: struct {
				name          string
				validatorName string
				sort          string
				desc          bool
				limit         uint64
				offset        uint64
			}{name: "pool1", validatorName: "val1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetValidator: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != dPool.ID, poolID is %s", poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition, epoch uint64) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						if condition.Sort.ValidatorDataSort != 1 {
							return nil, fmt.Errorf("condition.Sort.ValidatorDataSort != 1, condition.Sort.ValidatorDataSort is %d", condition.Sort.ValidatorDataSort)
						}
						if condition.Sort.Desc != true {
							return nil, fmt.Errorf("condition.Sort.Desc != true, condition.Sort.Desc is %v", false)
						}
						if condition.Condition.Name != "val1" {
							return nil, fmt.Errorf("condition.Condition.Name != val1, condition.Condition.Name is %s", condition.Condition.Name)
						}
						if condition.Condition.Pagination.Limit != 10 {
							return nil, fmt.Errorf("condition.Condition.Pagination.Limit != 10, but %d", condition.Condition.Pagination.Limit)
						}
						if condition.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("condition.Condition.Pagination.Offset != 10, but %d", condition.Condition.Pagination.Offset)
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string, epoch uint64) (*dmodels.ValidatorView, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"seventh": {
			Data: struct {
				name          string
				validatorName string
				sort          string
				desc          bool
				limit         uint64
				offset        uint64
			}{name: "pool1", validatorName: "val1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetValidatorDataCount: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != dPool.ID, poolID is %s", poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition, epoch uint64) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						if condition.Sort.ValidatorDataSort != 1 {
							return nil, fmt.Errorf("condition.Sort.ValidatorDataSort != 1, condition.Sort.ValidatorDataSort is %d", condition.Sort.ValidatorDataSort)
						}
						if condition.Sort.Desc != true {
							return nil, fmt.Errorf("condition.Sort.Desc != true, condition.Sort.Desc is %v", false)
						}
						if condition.Condition.Name != "val1" {
							return nil, fmt.Errorf("condition.Condition.Name != val1, condition.Condition.Name is %s", condition.Condition.Name)
						}
						if condition.Condition.Pagination.Limit != 10 {
							return nil, fmt.Errorf("condition.Condition.Pagination.Limit != 10, but %d", condition.Condition.Pagination.Limit)
						}
						if condition.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("condition.Condition.Pagination.Offset != 10, but %d", condition.Condition.Pagination.Offset)
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string, epoch uint64) (*dmodels.ValidatorView, error) {
						if validatorID != poolVD[0].ValidatorID {
							return nil, fmt.Errorf("validatorID != %s, validatorID is %s", poolVD[0].ValidatorID, validatorID)
						}
						return &dValView, nil
					},
					GetValidatorDataCountFunc: func(condition *postgres.PoolValidatorDataCondition, epoch uint64) (int64, error) {
						return 0, fmt.Errorf("some error")
					},
				},
			},
		},
	}

	for s, s2 := range data {
		t.Run(s, func(t *testing.T) {
			pv, count, err := s2.DAO.GetPoolValidators(s2.Data.name, s2.Data.validatorName, s2.Data.sort, s2.Data.desc, 1, s2.Data.limit, s2.Data.offset)
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			assert.Equal(t, uint64(len(s2.Result)), count)
			assert.Equal(t, uint64(len(pv)), count)
			for i, v := range pv {
				t.Run(fmt.Sprintf("governances[%d]", i), func(t *testing.T) {
					assert.DeepEqual(t, v, s2.Result[i])
				})
			}
		})
	}
}

func TestGetAllValidators(t *testing.T) {
	data := map[string]struct {
		DAO  services.Imp
		Data struct {
			validatorName string
			sort          string
			desc          bool
			limit         uint64
			offset        uint64
		}
		Result []*smodels.Validator
		Err    error
	}{
		"first": {
			Data: struct {
				validatorName string
				sort          string
				desc          bool
				limit         uint64
				offset        uint64
			}{validatorName: "val1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: []*smodels.Validator{
				{
					Name:             "val1",
					Image:            "img1",
					StakingAccounts:  500,
					NodePK:           "id1",
					APY:              decimal.Decimal{},
					VotePK:           "pk1",
					TotalActiveStake: sol.SOL{decimal.NewFromFloat(0.0000001)},
					Fee:              decimal.Decimal{},
					Score:            5698,
					SkippedSlots:     decimal.Decimal{},
					DataCenter:       "dc",
				},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetValidatorsFunc: func(condition *postgres.ValidatorCondition, epoch uint64) ([]*dmodels.ValidatorView, error) {
						if condition.Sort.ValidatorSort != 1 {
							return nil, fmt.Errorf("condition.Sort.ValidatorSort != 1, but %d", condition.Sort.ValidatorSort)
						}
						return []*dmodels.ValidatorView{&dValView}, nil
					},
					GetValidatorCountFunc: func(condition *postgres.ValidatorCondition, epoch uint64) (int64, error) {
						if condition.Condition.Name != "val1" {
							return 0, fmt.Errorf("condition.Condition.Name != val1, but %s", condition.Condition.Name)
						}
						return 1, nil
					},
				},
			},
		},
		"second": {
			Data: struct {
				validatorName string
				sort          string
				desc          bool
				limit         uint64
				offset        uint64
			}{validatorName: "val1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPoolValidatorData: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetValidatorsFunc: func(condition *postgres.ValidatorCondition, epoch uint64) ([]*dmodels.ValidatorView, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"third": {
			Data: struct {
				validatorName string
				sort          string
				desc          bool
				limit         uint64
				offset        uint64
			}{validatorName: "val1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: []*smodels.Validator{
				{
					Name:             "val1",
					Image:            "img1",
					StakingAccounts:  500,
					NodePK:           "id1",
					APY:              decimal.Decimal{},
					VotePK:           "pk1",
					TotalActiveStake: sol.SOL{decimal.NewFromFloat(0.0000001)},
					Fee:              decimal.Decimal{},
					Score:            5698,
					SkippedSlots:     decimal.Decimal{},
					DataCenter:       "dc",
				},
			},
			Err: fmt.Errorf("DAO.GetValidatorDataCount: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetValidatorsFunc: func(condition *postgres.ValidatorCondition, epoch uint64) ([]*dmodels.ValidatorView, error) {
						if condition.Sort.ValidatorSort != 1 {
							return nil, fmt.Errorf("condition.Sort.ValidatorSort != 1, but %d", condition.Sort.ValidatorSort)
						}
						return []*dmodels.ValidatorView{&dValView}, nil
					},
					GetValidatorCountFunc: func(condition *postgres.ValidatorCondition, epoch uint64) (int64, error) {
						return 0, fmt.Errorf("some error")
					},
				},
			},
		},
	}

	for s, s2 := range data {
		t.Run(s, func(t *testing.T) {
			gov, count, err := s2.DAO.GetAllValidators(s2.Data.validatorName, s2.Data.sort, s2.Data.desc, 1, []uint64{314}, s2.Data.limit, s2.Data.offset)
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
