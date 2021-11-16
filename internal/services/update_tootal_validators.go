package services

import (
	"context"
	"fmt"
	solana_sdk "github.com/everstake/solana-pools/pkg/extension/solana-sdk"
)

func (s Imp) UpdateValidators() error {

	ctx := context.Background()
	client := s.rpcClients["mainnet"]

	va, err := solana_sdk.GetVoteAccounts(client.RpcClient.Call(ctx, "getVoteAccounts"))
	if err != nil {
		return fmt.Errorf("UpdateValidators: %w", err)
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

	return nil
}
