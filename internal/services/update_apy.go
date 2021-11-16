package services

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
)

func (s Imp) UpdateAPY() error {

	ctx := context.Background()
	client := s.rpcClients["mainnet"]

	rate, err := client.RpcClient.GetInflationRate(ctx)
	if err != nil {
		return fmt.Errorf("UpdateAPY: %w", err)
	}

	s.cache.SetAPY(decimal.NewFromFloat(rate.Result.Total))

	return nil
}
