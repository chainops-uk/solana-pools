package services

import (
	"context"
	"fmt"
	"github.com/everstake/solana-pools/internal/services/smodels"
	solana_sdk "github.com/everstake/solana-pools/pkg/extension/solana-sdk"
	"github.com/shopspring/decimal"
	"time"
)

func (s Imp) UpdateNetworkData() error {

	ctx := context.Background()
	client := s.rpcClients["mainnet"]

	st, err := s.GetAvgSlotTimeMS()
	if err != nil {
		return fmt.Errorf("imp.GetAvgSlotTimeMS: %w", err)
	}

	ei, err := client.RpcClient.GetEpochInfo(ctx)
	if err != nil {
		return err
	}

	emptyS := ei.Result.SlotsInEpoch - ei.Result.SlotIndex

	progress := (float64(ei.Result.SlotIndex) / float64(ei.Result.SlotsInEpoch)) * 100

	sps := 1 / (st / 1000)

	s.Cache.SetCurrentEpochInfo(&smodels.EpochInfo{
		Epoch:        ei.Result.Epoch,
		SlotsInEpoch: ei.Result.SlotsInEpoch,
		SPS:          sps,
		EndEpoch:     time.Now().Add(time.Duration((float64(emptyS) / sps) * float64(time.Second))),
		Progress:     uint8(progress),
	})

	va, err := solana_sdk.GetVoteAccounts(client.RpcClient.Call(ctx, "getVoteAccounts"))
	if err != nil {
		return fmt.Errorf("GetVoteAccounts: %w", err)
	}

	var activeStake uint64
	for _, v := range va.Delinquent {
		activeStake += uint64(v.ActivatedStake)
	}
	for _, v := range va.Current {
		activeStake += uint64(v.ActivatedStake)
	}

	s.Cache.SetValidatorCount(int64(len(va.Current) + len(va.Delinquent)))
	s.Cache.SetActiveStake(activeStake)

	rate, err := client.RpcClient.GetInflationRate(ctx)
	if err != nil {
		return fmt.Errorf("GetInflationRate: %w", err)
	}

	sol, err := solana_sdk.GetSupply(client.RpcClient.Call(ctx, "getSupply"))
	if err != nil {
		return fmt.Errorf("GetSupply: %w", err)
	}

	apy := rate.Result.Total *
		(float64(sol.Value.Total) / float64(activeStake))
	APY := decimal.NewFromFloat(apy).Mul(decimal.NewFromInt(400).Div(decimal.NewFromFloat(st)))
	s.Cache.SetAPY(APY)

	return nil
}
