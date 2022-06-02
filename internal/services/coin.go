package services

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services/smodels"
	uuid "github.com/satori/go.uuid"
)

func (s Imp) GetPoolCoins(name string, sort string, desc bool, limit uint64, offset uint64) ([]*smodels.Coin, uint64, error) {

	pools, err := s.DAO.GetPools(&postgres.PoolCondition{
		Condition: &postgres.Condition{
			Network: postgres.MainNet,
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetPools: %w", err)
	}

	ids := make([]uuid.UUID, len(pools))
	for i, pool := range pools {
		ids[i] = pool.CoinID
	}

	coins, err := s.DAO.GetCoins(&postgres.CoinCondition{
		Condition: &postgres.Condition{
			IDs:        ids,
			Pagination: postgres.Pagination{Limit: limit, Offset: offset},
		},
		CoinSort: &postgres.CoinSort{
			Sort: postgres.SearchCoinSort(sort),
			Desc: desc,
		},
		Name: name,
	})

	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetCoins: %w", err)
	}

	scoins := make([]*smodels.Coin, len(coins))
	for i, coin := range coins {
		dEfI, err := s.DAO.GetDEFIs(&postgres.DeFiCondition{
			SaleCoinID: []uuid.UUID{coin.ID},
		})
		if err != nil {
			return nil, 0, fmt.Errorf("DAO.GetDEFIs: %w", err)
		}
		defi := make([]*smodels.DeFi, len(dEfI))
		for i2, d := range dEfI {
			lp, err := s.DAO.GetLiquidityPool(&postgres.Condition{IDs: []uuid.UUID{d.LiquidityPoolID}})
			if err != nil {
				return nil, 0, fmt.Errorf("DAO.GetLiquidityPool: %w", err)
			}
			coin, err := s.DAO.GetCoinByID(d.BuyCoinID)
			if err != nil {
				return nil, 0, fmt.Errorf("DAO.GetCoinByID: %w", err)
			}
			defi[i2] = (&smodels.DeFi{}).Set(d, (&smodels.Coin{}).Set(coin, nil), (&smodels.LiquidityPool{}).Set(lp))
		}

		scoins[i] = (&smodels.Coin{}).Set(coin, defi)
	}

	count, err := s.DAO.GetCoinsCount(&postgres.CoinCondition{
		Condition: &postgres.Condition{
			IDs: ids,
		},
		Name: name,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetCoinsCount: %w", err)
	}

	return scoins, uint64(count), nil
}

func (s Imp) GetCoins(name string, limit uint64, offset uint64) ([]*smodels.Coin, uint64, error) {
	coins, err := s.DAO.GetCoins(&postgres.CoinCondition{
		Condition: &postgres.Condition{
			Name:       name,
			Pagination: postgres.Pagination{Limit: limit, Offset: offset},
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetCoins: %w", err)
	}

	scoins := make([]*smodels.Coin, len(coins))
	for i, coin := range coins {
		scoins[i] = (&smodels.Coin{}).Set(coin, nil)
	}

	count, err := s.DAO.GetCoinsCount(&postgres.CoinCondition{
		Name: name,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("DAO.GetCoinsCount: %w", err)
	}

	return scoins, uint64(count), nil
}
