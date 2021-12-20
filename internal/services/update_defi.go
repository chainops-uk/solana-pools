package services

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"strings"
)

func (s Imp) UpdateDeFi() error {

	if err := updateRaydium(&s); err != nil {
		return fmt.Errorf("updateRaydium() %w", err)
	}

	return nil
}

func updateOrca(s *Imp) error {
	//pool, err := s.dao.GetLiquidityPool(&postgres.Condition{Names: []string{"Orca"}})
	//if err != nil {
	//	return err
	//}
	//if pool == nil {
	//	return nil
	//}
	//
	//pools, err := s.dao.GetPools(&postgres.Condition{Network: postgres.MainNet})
	//if err != nil {
	//	return fmt.Errorf("dao.GetPools: %w", err)
	//}
	//
	//ids := make([]uuid.UUID, len(pools))
	//for i, pool := range pools {
	//	ids[i] = pool.CoinID
	//}
	//
	//poolCoins, err := s.dao.GetCoins(&postgres.Condition{
	//	IDs: ids,
	//})
	//if err != nil {
	//	return fmt.Errorf("dao.GetCoins: %w", err)
	//}
	//
	//coins, err := s.dao.GetCoins(nil)
	//if err != nil {
	//	return err
	//}
	//
	//orca, err := s.orca.GetAllPools()
	//if err != nil {
	//	return err
	//}
	//
	////defis := make([]*dmodels.DEFI, 0)
	////for _, poolCoin := range poolCoins {
	////	for _, d := range coins {
	////		fi, ok := orca[fmt.Sprintf("%s/%s", poolCoin.Name, d.Name)]
	////		if !ok {
	////			continue
	////		}
	////		defis = append(defis, &dmodels.DEFI{
	////			LiquidityPoolID: pool.ID,
	////			SaleCoinID:      poolCoin.ID,
	////			BuyCoinID:       d.ID,
	////		})
	////	}
	////}
	//
	//
	//if err := s.dao.DeleteDeFis(&postgres.DeFiCondition{LiquidityPoolIDs: []uuid.UUID{pool.ID}}); err != nil {
	//	return err
	//}
	//
	//if err := s.dao.SaveDEFIs(defis...); err != nil {
	//	return err
	//}

	return nil
}

func updateRaydium(s *Imp) error {
	pool, err := s.dao.GetLiquidityPool(&postgres.Condition{Names: []string{"Raydium"}})
	if err != nil {
		return err
	}
	if pool == nil {
		return nil
	}

	pools, err := s.dao.GetPools(&postgres.Condition{Network: postgres.MainNet})
	if err != nil {
		return fmt.Errorf("dao.GetPools: %w", err)
	}

	ids := make([]uuid.UUID, len(pools))
	for i, pool := range pools {
		ids[i] = pool.CoinID
	}

	poolCoins, err := s.dao.GetCoins(&postgres.Condition{
		IDs: ids,
	})
	if err != nil {
		return fmt.Errorf("dao.GetCoins: %w", err)
	}

	coins, err := s.dao.GetCoins(nil)
	if err != nil {
		return err
	}

	raydium, err := s.raydium.GetPairs("")
	if err != nil {
		return err
	}

	defis := make([]*dmodels.DEFI, 0)
	for _, poolCoin := range poolCoins {
		for _, paris := range raydium {
			if strings.Contains(paris.PairID, fmt.Sprintf("%s-", poolCoin.Address)) {
				for _, d := range coins {
					address := d.Address
					if strings.Contains(address, "11111111111111111111111111111111") {
						address = "So" + address
					}
					if strings.Contains(paris.PairID, fmt.Sprintf("-%s", address)) {
						defis = append(defis, &dmodels.DEFI{
							LiquidityPoolID: pool.ID,
							SaleCoinID:      poolCoin.ID,
							BuyCoinID:       d.ID,
							Liquidity:       paris.Liquidity,
							APY:             decimal.NewFromFloat(paris.Apy),
						})
					}
				}
			}
		}

	}

	if err := s.dao.DeleteDeFis(&postgres.DeFiCondition{LiquidityPoolIDs: []uuid.UUID{pool.ID}}); err != nil {
		return err
	}

	if err := s.dao.SaveDEFIs(defis...); err != nil {
		return err
	}

	return nil
}
