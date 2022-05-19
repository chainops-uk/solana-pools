package services

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services/smodels"
	uuid "github.com/satori/go.uuid"
)

func (s Imp) GetPoolValidators(name string, validatorName string, sort string, desc bool, limit uint64, offset uint64) ([]*smodels.PoolValidatorData, uint64, error) {
	pool, err := s.dao.GetPool(name)
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetPool: %w", err)
	}
	if pool == nil {
		return nil, 0, fmt.Errorf("dao.GetPool(%s): %w", name, postgres.ErrorRecordNotFounded)
	}

	poolData, err := s.dao.GetLastPoolData(pool.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetLastPoolData: %w", err)
	}

	pvd, err := s.dao.GetPoolValidatorData(&postgres.PoolValidatorDataCondition{
		PoolDataIDs: []uuid.UUID{poolData.ID},
		Sort: &postgres.ValidatorDataSort{
			ValidatorDataSort: postgres.SearchValidatorDataSort(sort),
			Desc:              desc,
		},
		Condition: &postgres.Condition{
			Name: validatorName,
			Pagination: postgres.Pagination{
				Limit:  limit,
				Offset: offset,
			},
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetPoolValidatorData: %w", err)
	}

	arr := make([]*smodels.PoolValidatorData, len(pvd))
	for i, data := range pvd {
		val, err := s.dao.GetValidator(data.ValidatorID)
		if err != nil {
			return nil, 0, fmt.Errorf("dao.GetValidator: %w", err)
		}
		arr[i] = (&smodels.PoolValidatorData{}).Set(data.ActiveStake, val)
	}

	count, err := s.dao.GetValidatorDataCount(&postgres.PoolValidatorDataCondition{
		PoolDataIDs: []uuid.UUID{
			poolData.ID,
		},
		Condition: &postgres.Condition{
			Name: validatorName,
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetValidatorDataCount: %w", err)
	}

	return arr, uint64(count), nil
}

func (s Imp) GetAllValidators(validatorName string, sort string, desc bool, limit uint64, offset uint64) ([]*smodels.Validator, uint64, error) {
	pvd, err := s.dao.GetValidators(&postgres.ValidatorCondition{
		Sort: &postgres.ValidatorSort{
			ValidatorSort: postgres.SearchValidatorSort(sort),
			Desc:          desc,
		},
		Condition: &postgres.Condition{
			Name: validatorName,
			Pagination: postgres.Pagination{
				Limit:  limit,
				Offset: offset,
			},
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetPoolValidatorData: %w", err)
	}

	arr := make([]*smodels.Validator, len(pvd))
	for i, data := range pvd {
		arr[i] = (&smodels.Validator{}).Set(data)
	}

	count, err := s.dao.GetValidatorCount(&postgres.ValidatorCondition{
		Condition: &postgres.Condition{
			Name: validatorName,
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetValidatorDataCount: %w", err)
	}

	return arr, uint64(count), nil
}
