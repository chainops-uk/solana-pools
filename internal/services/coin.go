package services

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services/smodels"
	uuid "github.com/satori/go.uuid"
)

func (s Imp) GetPoolCoins(name string, limit uint64, offset uint64) ([]*smodels.Coin, uint64, error) {

	pools, err := s.dao.GetPools(&postgres.Condition{Network: postgres.MainNet})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetPools: %w", err)
	}

	ids := make([]uuid.UUID, len(pools))
	for i, pool := range pools {
		ids[i] = pool.CoinID
	}

	coins, err := s.dao.GetCoins(&postgres.Condition{
		IDs:        ids,
		Name:       name,
		Pagination: postgres.Pagination{Limit: limit, Offset: offset},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetCoins: %w", err)
	}

	scoins := make([]*smodels.Coin, len(coins))
	for i, coin := range coins {
		scoins[i] = (&smodels.Coin{}).Set(coin)
	}

	count, err := s.dao.GetCoinsCount(&postgres.Condition{
		IDs:  ids,
		Name: name,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetCoinsCount: %w", err)
	}

	return scoins, uint64(count), nil
}

func (s Imp) GetCoins(name string, limit uint64, offset uint64) ([]*smodels.Coin, uint64, error) {
	coins, err := s.dao.GetCoins(&postgres.Condition{
		Name:       name,
		Pagination: postgres.Pagination{Limit: limit, Offset: offset},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetCoins: %w", err)
	}

	scoins := make([]*smodels.Coin, len(coins))
	for i, coin := range coins {
		scoins[i] = (&smodels.Coin{}).Set(coin)
	}

	count, err := s.dao.GetCoinsCount(&postgres.Condition{
		Name: name,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("dao.GetCoinsCount: %w", err)
	}

	return scoins, uint64(count), nil
}
