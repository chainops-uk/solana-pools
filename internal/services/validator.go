package services

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services/smodels"
	uuid "github.com/satori/go.uuid"
)

func (s Imp) GetPoolValidators(name string, validatorName string, sort string, desc bool, epoch uint64, limit uint64, offset uint64) ([]*smodels.PoolValidatorData, uint64, error) {
	pool, err := s.DAO.GetPool(name)
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetPool: %w", err)
	}

	if pool == nil {
		return nil, 0, fmt.Errorf("DAO.GetPool(%s): %w", name, postgres.ErrorRecordNotFounded)
	}

	poolData, err := s.DAO.GetLastPoolData(pool.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetLastPoolData: %w", err)
	}

	pvd, err := s.DAO.GetPoolValidatorData(&postgres.PoolValidatorDataCondition{
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
	}, epoch)
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetPoolValidatorData: %w", err)
	}

	arr := make([]*smodels.PoolValidatorData, len(pvd))
	for i, data := range pvd {
		val, err := s.DAO.GetValidator(data.ValidatorID, epoch)
		if err != nil {
			return nil, 0, fmt.Errorf("DAO.GetValidator: %w", err)
		}
		arr[i] = (&smodels.PoolValidatorData{}).Set(data.ActiveStake, val)
	}

	count, err := s.DAO.GetValidatorDataCount(&postgres.PoolValidatorDataCondition{
		PoolDataIDs: []uuid.UUID{
			poolData.ID,
		},
		Sort: &postgres.ValidatorDataSort{
			ValidatorDataSort: postgres.SearchValidatorDataSort(sort),
			Desc:              desc,
		},
		Condition: &postgres.Condition{
			Name: validatorName,
		},
	}, epoch)
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetValidatorDataCount: %w", err)
	}

	return arr, uint64(count), nil
}

func (s Imp) GetAllValidators(validatorName string, sort string, desc bool, epoch uint64, epochs []uint64, limit uint64, offset uint64) ([]*smodels.Validator, uint64, error) {
	pvd, err := s.DAO.GetValidators(&postgres.ValidatorCondition{
		Epochs: epochs,
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
	}, epoch)
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetPoolValidatorData: %w", err)
	}

	arr := make([]*smodels.Validator, len(pvd))
	for i, data := range pvd {
		arr[i] = (&smodels.Validator{}).Set(data)
	}

	count, err := s.DAO.GetValidatorCount(&postgres.ValidatorCondition{
		Epochs: epochs,
		Condition: &postgres.Condition{
			Name: validatorName,
		},
	}, epoch)
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetValidatorDataCount: %w", err)
	}

	return arr, uint64(count), nil
}
