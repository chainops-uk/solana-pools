package services

import "fmt"

func (s Imp) UpdateCoins() error {
	coins, err := s.dao.GetCoins(nil)
	if err != nil {
		return fmt.Errorf("dao.GetCoins: %w", err)
	}

	for _, coin := range coins {
		if coin.Address == "00000000000000000000000000000000" {
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

	if err := s.dao.SaveCoin(coins...); err != nil {
		return fmt.Errorf("dao.SaveCoin: %w", err)
	}

	return nil
}
