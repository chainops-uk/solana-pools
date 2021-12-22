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

	ei1, err := client.RpcClient.GetEpochInfo(ctx)
	if err != nil {
		return err
	}

	t1 := time.Now()

	<-time.After(time.Minute * 1)

	ei2, err := client.RpcClient.GetEpochInfo(ctx)
	if err != nil {
		return err
	}

	t2 := time.Now()

	if ei1.Result.Epoch != ei2.Result.Epoch {
		return err
	}

	sps := float64(ei2.Result.SlotIndex-ei1.Result.SlotIndex) / t2.Sub(t1).Seconds()

	emptyS := ei2.Result.SlotsInEpoch - ei2.Result.SlotIndex

	progress := (float64(ei2.Result.SlotIndex) / float64(ei2.Result.SlotsInEpoch)) * 100

	if sps == 0 {
		return err
	}

	s.cache.SetCurrentEpochInfo(&smodels.EpochInfo{
		Epoch:        ei2.Result.Epoch,
		SlotsInEpoch: ei2.Result.SlotsInEpoch,
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

	s.cache.SetValidatorCount(int64(len(va.Current) + len(va.Delinquent)))
	s.cache.SetActiveStake(activeStake)

	rate, err := client.RpcClient.GetInflationRate(ctx)
	if err != nil {
		return fmt.Errorf("GetInflationRate: %w", err)
	}

	sol, err := solana_sdk.GetSupply(client.RpcClient.Call(ctx, "getSupply"))
	if err != nil {
		return fmt.Errorf("GetSupply: %w", err)
	}

	st, err := s.GetAvgSlotTimeMS()
	if err != nil {
		return fmt.Errorf("imp.GetAvgSlotTimeMS: %w", err)
	}

	apy := rate.Result.Total *
		(float64(sol.Value.Total) / float64(activeStake))
	APY := decimal.NewFromFloat(apy).Mul(decimal.NewFromInt(400).Div(decimal.NewFromFloat(st)))
	s.cache.SetAPY(APY)

	return nil
}
