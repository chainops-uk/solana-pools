package services

import "fmt"

func (s Imp) UpdateCoins() error {
	coins, err := s.DAO.GetCoins(nil)
	if err != nil {
		return fmt.Errorf("DAO.GetCoins: %w", err)
	}

	for _, coin := range coins {
		if coin.GeckoKey == "null" {
			continue
		}
		c, err := s.coinGecko.CoinsID(coin.GeckoKey,
			false,
			false,
			true,
			false,
			false,
			false)
		if err != nil {
			return fmt.Errorf("coinGecko.CoinsID(%s): %w", coin.GeckoKey, err)
		}

		coin.USD = c.MarketData.CurrentPrice["usd"]

		coin.LargeImage = c.Image.Large
		coin.ThumbImage = c.Image.Thumb
		coin.SmallImage = c.Image.Small
	}

	if err := s.DAO.SaveCoin(coins...); err != nil {
		return fmt.Errorf("DAO.SaveCoin: %w", err)
	}

	return nil
}
