package services

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

func (s Imp) UpdatePrice() error {
	coin, err := s.coinGecko.CoinsID("solana",
		false,
		false,
		true,
		false,
		false,
		false)
	if err != nil {
		return fmt.Errorf("UpdatePrice: %w", err)
	}

	usd, ok := coin.MarketData.CurrentPrice["usd"]
	if !ok {
		return fmt.Errorf("UpdatePrice: %w", errors.New("usd price not found"))
	}

	s.Cache.SetPrice(decimal.NewFromFloat(usd))

	return nil
}
