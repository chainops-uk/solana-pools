package services

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services/smodels"
)

func (s Imp) GetLiquidityPools(name string, limit uint64, offset uint64) ([]*smodels.LiquidityPool, uint64, error) {
	gov, err := s.DAO.GetLiquidityPools(&postgres.Condition{
		Name:       name,
		Pagination: postgres.Pagination{Limit: limit, Offset: offset},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetLiquidityPools: %w", err)
	}

	arr := make([]*smodels.LiquidityPool, len(gov))
	for i, g := range gov {
		arr[i] = (&smodels.LiquidityPool{}).Set(g)
	}

	count, err := s.DAO.GetLiquidityPoolsCount(&postgres.Condition{Name: name})
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetLiquidityPoolsCount: %w", err)
	}

	return arr, uint64(count), nil
}
