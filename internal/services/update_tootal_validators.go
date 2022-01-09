package services

import (
	"context"
	"fmt"
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	solana_sdk "github.com/everstake/solana-pools/pkg/extension/solana-sdk"
	"github.com/everstake/solana-pools/pkg/validatorsapp"
	"github.com/shopspring/decimal"
	"math"
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
			VotePK:          v.VotePubKey,
			APY:             apy,
			StakingAccounts: stakingAccounts,
			ActiveStake:     uint64(v.ActivatedStake),
			Fee:             decimal.NewFromFloat(float64(vInfo.Commission) / 100.0),
			Score:           vInfo.TotalScore,
			SkippedSlots:    skippedSlots,
			DataCenter:      vInfo.DataCenterKey,
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

		apy = apy.Mul(decimal.NewFromFloat(correlation))

		validators = append(validators, &dmodels.Validator{
			ID:              v.NodePubKey,
			Name:            vInfo.Name,
			Image:           vInfo.AvatarURL,
			Delinquent:      true,
			VotePK:          v.VotePubKey,
			APY:             apy,
			StakingAccounts: stakingAccounts,
			ActiveStake:     uint64(v.ActivatedStake),
			Fee:             decimal.NewFromFloat(float64(vInfo.Commission) / 100.0),
			Score:           vInfo.TotalScore,
			SkippedSlots:    skippedSlots,
			DataCenter:      vInfo.DataCenterKey,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		})
	}

	step := 100

	n := int(math.Ceil(float64(len(validators)) / float64(step)))
	offset := 0
	for i := 0; i < n; i++ {
		if offset+step > len(validators) {
			if err := s.dao.UpdateValidators(validators[offset:]...); err != nil {
				return fmt.Errorf("dao.UpdateValidators: %w", err)
			}
		} else {
			if err := s.dao.UpdateValidators(validators[offset : offset+step]...); err != nil {
				return fmt.Errorf("dao.UpdateValidators: %w", err)
			}
		}

		offset += step
	}

	return nil
}
