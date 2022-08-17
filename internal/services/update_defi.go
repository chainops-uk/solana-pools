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
	if err := updateOrca(&s); err != nil {
		return fmt.Errorf("updateOrca() %w", err)
	}
	if err := updateAtrix(&s); err != nil {
		return fmt.Errorf("updateAtrix() %w", err)
	}
	if err := updateSaber(&s); err != nil {
		return fmt.Errorf("updateSaber() %w", err)
	}

	return nil
}

func updateOrca(s *Imp) error {
	pool, err := s.DAO.GetLiquidityPool(&postgres.Condition{Names: []string{"Orca"}})
	if err != nil {
		return err
	}
	if pool == nil {
		return nil
	}

	pools, err := s.DAO.GetPools(&postgres.PoolCondition{Condition: &postgres.Condition{Network: postgres.MainNet}})
	if err != nil {
		return fmt.Errorf("DAO.GetPools: %w", err)
	}

	ids := make([]uuid.UUID, len(pools))
	for i, pool := range pools {
		ids[i] = pool.CoinID
	}

	poolCoins, err := s.DAO.GetCoins(&postgres.CoinCondition{
		Condition: &postgres.Condition{
			IDs: ids,
		},
	})
	if err != nil {
		return fmt.Errorf("DAO.GetCoins: %w", err)
	}

	coins, err := s.DAO.GetCoins(nil)
	if err != nil {
		return err
	}

	orca, err := s.orca.GetPools()
	if err != nil {
		return err
	}

	defis := make([]*dmodels.DEFI, 0)
	for _, poolCoin := range poolCoins {
		for _, o := range orca {
			if strings.Contains(o.Name, fmt.Sprintf("%s/", poolCoin.Name)) {
				for _, d := range coins {
					if strings.Contains(o.Name, fmt.Sprintf("/%s", d.Name)) {
						apy := decimal.NewFromInt(0)
						if o.Apy24H != nil {
							apy = decimal.NewFromFloat(*o.Apy24H)
						}
						defis = append(defis, &dmodels.DEFI{
							LiquidityPoolID: pool.ID,
							SaleCoinID:      poolCoin.ID,
							BuyCoinID:       d.ID,
							Liquidity:       o.Liquidity,
							APY:             apy,
						})
					}
				}
			}
		}

	}

	if err := s.DAO.DeleteDeFis(&postgres.DeFiCondition{LiquidityPoolIDs: []uuid.UUID{pool.ID}}); err != nil {
		return err
	}

	if err := s.DAO.SaveDEFIs(defis...); err != nil {
		return err
	}

	return nil
}

func updateRaydium(s *Imp) error {
	pool, err := s.DAO.GetLiquidityPool(&postgres.Condition{Names: []string{"Raydium"}})
	if err != nil {
		return err
	}
	if pool == nil {
		return nil
	}

	pools, err := s.DAO.GetPools(&postgres.PoolCondition{Condition: &postgres.Condition{Network: postgres.MainNet}})
	if err != nil {
		return fmt.Errorf("DAO.GetPools: %w", err)
	}

	ids := make([]uuid.UUID, len(pools))
	for i, pool := range pools {
		ids[i] = pool.CoinID
	}

	poolCoins, err := s.DAO.GetCoins(&postgres.CoinCondition{
		Condition: &postgres.Condition{
			IDs: ids,
		},
	})
	if err != nil {
		return fmt.Errorf("DAO.GetCoins: %w", err)
	}

	coins, err := s.DAO.GetCoins(nil)
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
						address = "So11111111111111111111111111111111"
					}
					if strings.Contains(paris.PairID, fmt.Sprintf("-%s", address)) {
						defis = append(defis, &dmodels.DEFI{
							LiquidityPoolID: pool.ID,
							SaleCoinID:      poolCoin.ID,
							BuyCoinID:       d.ID,
							Liquidity:       paris.Liquidity,
							APY:             decimal.NewFromFloat(paris.Apy).Div(decimal.NewFromInt(100)),
						})
					}
				}
			}
		}

	}

	if err := s.DAO.DeleteDeFis(&postgres.DeFiCondition{LiquidityPoolIDs: []uuid.UUID{pool.ID}}); err != nil {
		return err
	}

	if err := s.DAO.SaveDEFIs(defis...); err != nil {
		return err
	}

	return nil
}

func updateAtrix(s *Imp) error {
	pool, err := s.DAO.GetLiquidityPool(&postgres.Condition{Names: []string{"Atrix"}})
	if err != nil {
		return err
	}
	if pool == nil {
		return nil
	}

	pools, err := s.DAO.GetPools(&postgres.PoolCondition{Condition: &postgres.Condition{Network: postgres.MainNet}})
	if err != nil {
		return fmt.Errorf("DAO.GetPools: %w", err)
	}

	ids := make([]uuid.UUID, len(pools))
	for i, pool := range pools {
		ids[i] = pool.CoinID
	}

	poolCoins, err := s.DAO.GetCoins(&postgres.CoinCondition{
		Condition: &postgres.Condition{
			IDs: ids,
		},
	})
	if err != nil {
		return fmt.Errorf("DAO.GetCoins: %w", err)
	}

	coins, err := s.DAO.GetCoins(nil)
	if err != nil {
		return err
	}

	atrix, err := s.atrix.GetTVL()
	if err != nil {
		return err
	}

	defis := make([]*dmodels.DEFI, 0)
	for _, poolCoin := range poolCoins {
		for _, v := range atrix.Pools {
			if v.CoinMint == poolCoin.Address {
				for _, d := range coins {
					if d.Address == v.PCMint {
						defis = append(defis, &dmodels.DEFI{
							LiquidityPoolID: pool.ID,
							SaleCoinID:      poolCoin.ID,
							BuyCoinID:       d.ID,
							Liquidity:       v.Tvl,
							APY:             decimal.NewFromInt(0),
						})
					}
				}
			}
		}

	}

	if err := s.DAO.DeleteDeFis(&postgres.DeFiCondition{LiquidityPoolIDs: []uuid.UUID{pool.ID}}); err != nil {
		return err
	}

	if err := s.DAO.SaveDEFIs(defis...); err != nil {
		return err
	}

	return nil
}

func updateSaber(s *Imp) error {
	pool, err := s.DAO.GetLiquidityPool(&postgres.Condition{Names: []string{"Saber"}})
	if err != nil {
		return err
	}
	if pool == nil {
		return nil
	}

	pools, err := s.DAO.GetPools(&postgres.PoolCondition{Condition: &postgres.Condition{Network: postgres.MainNet}})
	if err != nil {
		return fmt.Errorf("DAO.GetPools: %w", err)
	}

	ids := make([]uuid.UUID, len(pools))
	for i, pool := range pools {
		ids[i] = pool.CoinID
	}

	poolCoins, err := s.DAO.GetCoins(&postgres.CoinCondition{
		Condition: &postgres.Condition{
			IDs: ids,
		},
	})
	if err != nil {
		return fmt.Errorf("DAO.GetCoins: %w", err)
	}

	coins, err := s.DAO.GetCoins(nil)
	if err != nil {
		return err
	}

	saber, err := s.saber.GetPools()
	if err != nil {
		return err
	}

	defis := make([]*dmodels.DEFI, 0)
	for _, poolCoin := range poolCoins {
		for _, v := range saber {
			if v.Coin.Address == poolCoin.Address {
				for _, d := range coins {
					if d.Address == v.PC.Address {
						defis = append(defis, &dmodels.DEFI{
							LiquidityPoolID: pool.ID,
							SaleCoinID:      poolCoin.ID,
							BuyCoinID:       d.ID,
							Liquidity:       v.Stats.TvlCoin*poolCoin.USD + v.Stats.TvlPC*d.USD,
							APY: decimal.NewFromFloat(v.Stats.Vol24H * d.USD).
								Mul(decimal.NewFromFloat(0.0004)).
								Div(decimal.NewFromFloat(v.Stats.TvlCoin*poolCoin.USD + v.Stats.TvlPC*d.USD)).
								Mul(decimal.NewFromInt(365)),
						})
					}
				}
			}
		}

	}

	if err := s.DAO.DeleteDeFis(&postgres.DeFiCondition{LiquidityPoolIDs: []uuid.UUID{pool.ID}}); err != nil {
		return err
	}

	if err := s.DAO.SaveDEFIs(defis...); err != nil {
		return err
	}

	return nil
}
