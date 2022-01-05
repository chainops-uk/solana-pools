package services

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services/smodels"
)

func (s Imp) GetGovernance(name string, sort string, desc bool, limit uint64, offset uint64) ([]*smodels.Governance, uint64, error) {
	gov, err := s.dao.GetGovernance(&postgres.GovernanceCondition{
		Condition: &postgres.Condition{
			Name:       name,
			Pagination: postgres.Pagination{Limit: limit, Offset: offset},
		},
		Sort: &postgres.GovernanceSort{
			Sort: postgres.SearchGovernanceSort(sort),
			Desc: desc,
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetCoins: %w", err)
	}

	sgov := make([]*smodels.Governance, len(gov))
	for i, g := range gov {
		sgov[i] = (&smodels.Governance{}).Set(g)
	}

	count, err := s.dao.GetGovernanceCount(&postgres.GovernanceCondition{
		Condition: &postgres.Condition{
			Name: name,
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetCoinsCount: %w", err)
	}

	return sgov, uint64(count), nil
}
