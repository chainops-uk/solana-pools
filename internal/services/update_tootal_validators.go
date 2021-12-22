package services

import (
	"context"
	"fmt"
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	solana_sdk "github.com/everstake/solana-pools/pkg/extension/solana-sdk"
	"github.com/everstake/solana-pools/pkg/validatorsapp"
	"github.com/shopspring/decimal"
	"time"
)

func (s Imp) UpdateValidators() error {

	ctx := context.Background()
	client := s.rpcClients["mainnet"]

	st, err := s.GetAvgSlotTimeMS()
	if err != nil {
		return fmt.Errorf("imp.GetAvgSlotTimeMS: %w", err)
	}
	correlation := 400 / st

	va, err := solana_sdk.GetVoteAccounts(client.RpcClient.Call(ctx, "getVoteAccounts"))
	if err != nil {
		return fmt.Errorf("UpdateValidators: %w", err)
	}

	validators := make([]*dmodels.Validator, 0, len(va.Current)+len(va.Delinquent))

	for _, v := range va.Current {
		var vInfo validatorsapp.ValidatorAppInfo
		err = rep(func() error {
			vInfo, err = s.validatorsApp.GetValidatorInfo("mainnet", v.NodePubKey)
			return err
		}, 10, time.Minute*1)
		if err != nil {
			return fmt.Errorf("validatorsApp.GetValidatorInfo(%s): %w", v.NodePubKey, err)
		}
		skippedSlots, _ := decimal.NewFromString(vInfo.SkippedSlotPercent)
		apy, stakingAccounts, err := getAPY(client, ctx, solana.MustPublicKeyFromBase58(v.VotePubKey), EpochsPerYear)
		if err != nil {
			return fmt.Errorf("getAPY: %w", err)
		}

		apy = apy.Mul(decimal.NewFromFloat(correlation))

		validators = append(validators, &dmodels.Validator{
			ID:              v.NodePubKey,
			Name:            vInfo.Name,
			Image:           vInfo.AvatarURL,
			Delinquent:      false,
			Network:         "mainnet",
			VotePK:          v.VotePubKey,
			APY:             apy,
			StakingAccounts: stakingAccounts,
			ActiveStake:     uint64(v.ActivatedStake),
			Fee:             decimal.NewFromFloat(float64(vInfo.Commission) / 100.0),
			Score:           vInfo.TotalScore,
			SkippedSlots:    skippedSlots,
			DataCenter:      vInfo.DataCenterHost,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		})
	}
	for _, v := range va.Delinquent {
		var vInfo validatorsapp.ValidatorAppInfo
		err = rep(func() error {
			vInfo, err = s.validatorsApp.GetValidatorInfo("mainnet", v.NodePubKey)
			return err
		}, 10, time.Minute*1)
		if err != nil {
			return fmt.Errorf("validatorsApp.GetValidatorInfo(%s): %w", v.NodePubKey, err)
		}
		skippedSlots, _ := decimal.NewFromString(vInfo.SkippedSlotPercent)
		apy, stakingAccounts, err := getAPY(client, ctx, solana.MustPublicKeyFromBase58(v.VotePubKey), EpochsPerYear)
		if err != nil {
			return fmt.Errorf("getAPY: %w", err)
		}

		validators = append(validators, &dmodels.Validator{
			ID:              v.NodePubKey,
			Name:            vInfo.Name,
			Delinquent:      true,
			Network:         "mainnet",
			VotePK:          v.VotePubKey,
			APY:             apy,
			StakingAccounts: stakingAccounts,
			ActiveStake:     uint64(v.ActivatedStake),
			Fee:             decimal.NewFromFloat(float64(vInfo.Commission) / 100.0),
			Score:           vInfo.TotalScore,
			SkippedSlots:    skippedSlots,
			DataCenter:      vInfo.DataCenterHost,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		})
	}

	if err := s.dao.UpdateValidators(validators...); err != nil {
		return fmt.Errorf("dao.UpdateValidators: %w", err)
	}

	return nil
}

func (s Imp) UpdateTestNetValidators() error {

	ctx := context.Background()
	client := s.rpcClients["testnet"]

rep:
	ei1, err := client.RpcClient.GetEpochInfo(ctx)
	if err != nil {
		return fmt.Errorf("GetEpochInfo: %w", err)
	}

	t1 := time.Now()

	<-time.After(time.Minute * 1)

	ei2, err := client.RpcClient.GetEpochInfo(ctx)
	if err != nil {
		return fmt.Errorf("GetEpochInfo: %w", err)
	}

	t2 := time.Now()

	if ei1.Result.Epoch != ei2.Result.Epoch {
		goto rep
	}

	sps := float64(ei2.Result.SlotIndex-ei1.Result.SlotIndex) / t2.Sub(t1).Seconds()

	epochTime := float64(ei2.Result.SlotsInEpoch) / sps

	epochInYear := 31557600 / epochTime

	va, err := solana_sdk.GetVoteAccounts(client.RpcClient.Call(ctx, "getVoteAccounts"))
	if err != nil {
		return fmt.Errorf("UpdateValidators: %w", err)
	}

	ii := 0
	validators := make([]*dmodels.Validator, len(va.Current)+len(va.Delinquent))
	for _, v := range va.Current {
		var vInfo validatorsapp.ValidatorAppInfo
		err = rep(func() error {
			vInfo, err = s.validatorsApp.GetValidatorInfo("testnet", v.NodePubKey)
			return err
		}, 10, time.Minute*1)
		if err != nil {
			return fmt.Errorf("validatorsApp.GetValidatorInfo(%s): %w", v.NodePubKey, err)
		}
		skippedSlots, _ := decimal.NewFromString(vInfo.SkippedSlotPercent)
		apy, stakingAccounts, err := getAPY(client, ctx, solana.MustPublicKeyFromBase58(v.VotePubKey), epochInYear)
		if err != nil {
			return fmt.Errorf("getAPY: %w", err)
		}

		validators[ii] = &dmodels.Validator{
			ID:              v.NodePubKey,
			Name:            vInfo.Name,
			Delinquent:      false,
			Network:         "testnet",
			VotePK:          v.VotePubKey,
			APY:             apy,
			StakingAccounts: stakingAccounts,
			ActiveStake:     uint64(v.ActivatedStake),
			Fee:             decimal.NewFromFloat(float64(vInfo.Commission) / 100.0),
			Score:           vInfo.TotalScore,
			SkippedSlots:    skippedSlots,
			DataCenter:      vInfo.DataCenterHost,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		ii++
	}
	for _, v := range va.Delinquent {
		var vInfo validatorsapp.ValidatorAppInfo
		err = rep(func() error {
			vInfo, err = s.validatorsApp.GetValidatorInfo("mainnet", v.NodePubKey)
			return err
		}, 10, time.Minute*1)
		if err != nil {
			return fmt.Errorf("validatorsApp.GetValidatorInfo(%s): %w", v.NodePubKey, err)
		}
		skippedSlots, _ := decimal.NewFromString(vInfo.SkippedSlotPercent)
		apy, stakingAccounts, err := getAPY(client, ctx, solana.MustPublicKeyFromBase58(v.VotePubKey), epochInYear)
		if err != nil {
			return fmt.Errorf("getAPY: %w", err)
		}

		validators[ii] = &dmodels.Validator{
			ID:              v.NodePubKey,
			Name:            vInfo.Name,
			Delinquent:      true,
			Network:         "mainnet",
			VotePK:          v.VotePubKey,
			APY:             apy,
			StakingAccounts: stakingAccounts,
			ActiveStake:     uint64(v.ActivatedStake),
			Fee:             decimal.NewFromFloat(float64(vInfo.Commission) / 100.0),
			Score:           vInfo.TotalScore,
			SkippedSlots:    skippedSlots,
			DataCenter:      vInfo.DataCenterHost,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		ii++
	}

	if err := s.dao.UpdateValidators(validators...); err != nil {
		return fmt.Errorf("dao.UpdateValidators: %w", err)
	}

	return nil
}
