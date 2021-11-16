package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	solana_sdk "github.com/everstake/solana-pools/pkg/extension/solana-sdk"
	"github.com/everstake/solana-pools/pkg/pools"
	"github.com/everstake/solana-pools/pkg/pools/types"
	"github.com/everstake/solana-pools/pkg/validatorsapp"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"math"
	"time"
)

func (s Imp) UpdatePools() error {
	dPools, err := s.dao.GetPools(nil)
	if err != nil {
		return fmt.Errorf("dao.GetPools: %s", err.Error())
	}
	var success, fail uint64
	start := time.Now()
	for _, p := range dPools {
		if !p.Active {
			continue
		}
		if err := s.updatePool(p); err != nil {
			s.log.Error(
				"Update Pools",
				zap.String("pool_name", p.Name),
				zap.String("pool_address", p.Address),
				zap.String("network", p.Network),
				zap.Error(err),
			)
			fail++
		} else {
			success++
		}
	}
	s.log.Debug(
		"Pools Updated",
		zap.Uint64("success", success),
		zap.Uint64("failed", fail),
		zap.Duration("duration", time.Now().Sub(start)),
	)
	return nil
}

func (s Imp) updatePool(dPool dmodels.Pool) error {
	ctx := context.Background()
	net := config.Network(dPool.Network)
	rpcCli, ok := s.rpcClients[net]
	if !ok {
		return fmt.Errorf("rpc client for %s network not found", dPool.Network)
	}
	poolFactory := pools.NewFactory(rpcCli)
	pool, err := poolFactory.GetPool(dPool.Name)
	if err != nil {
		return fmt.Errorf("poolFactory.GetPool: %s", err.Error())
	}
	data, err := pool.GetData(dPool.Address)
	if err != nil {
		return fmt.Errorf("pool.GetData: %s", err.Error())
	}

rep:
	ei1, err := rpcCli.RpcClient.GetEpochInfo(ctx)
	if err != nil {
		return fmt.Errorf("GetEpochInfo: %w", err)
	}

	t1 := time.Now()

	<-time.After(time.Minute * 1)

	ei2, err := rpcCli.RpcClient.GetEpochInfo(ctx)
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

	dmodel := &dmodels.PoolData{
		ID:            uuid.NewV4(),
		PoolID:        dPool.ID,
		ActiveStake:   lampToSol(data.SolanaStake),
		TokensSupply:  decimal.New(int64(data.TokenSupply), -9),
		DepossitFee:   decimal.NewFromFloat(data.DepositFee).Truncate(-2),
		WithdrawalFee: decimal.NewFromFloat(data.WithdrawalFee).Truncate(-2),
		RewardsFee:    decimal.NewFromFloat(data.RewardsFee).Truncate(-2),
		UpdatedAt:     time.Now(),
		CreatedAt:     time.Now(),
	}

	var avgSkippedSlots decimal.Decimal
	var avgScore int64
	var delinquent int64
	validators := make([]*dmodels.Validator, len(data.Validators))
	for i, v := range data.Validators {
		if v.NodePK == types.EmptyAddress {
			v.NodePK, err = getNodeAddress(rpcCli, ctx, v.VotePK)
			if err != nil {
				return fmt.Errorf("getNodeAddress(%s): %s", v.VotePK, err.Error())
			}
		}
		var vInfo validatorsapp.ValidatorAppInfo
		err = rep(func() error {
			vInfo, err = s.validatorsApp.GetValidatorInfo(dPool.Network, v.NodePK.String())
			return err
		}, 10, time.Minute*1)
		if err != nil {
			return fmt.Errorf("validatorsApp.GetValidatorInfo(%s): %s", v.NodePK, err.Error())
		}
		skippedSlots, _ := decimal.NewFromString(vInfo.SkippedSlotPercent)
		validators[i] = &dmodels.Validator{
			PoolDataID:   dmodel.ID,
			VotePK:       v.VotePK.String(),
			NodePK:       v.NodePK.String(),
			ActiveStake:  lampToSol(v.ActiveStake),
			Fee:          decimal.New(vInfo.Commission, 0),
			Score:        vInfo.TotalScore,
			SkippedSlots: skippedSlots,
			DataCenter:   vInfo.DataCenterHost,
		}

		validators[i].APY, validators[i].StakeAccounts, err = getAPY(rpcCli, ctx, v.VotePK, epochInYear)
		if err != nil {
			return fmt.Errorf("getAPY: %w", err)
		}

		if vInfo.Delinquent {
			delinquent++
		}
		avgSkippedSlots = avgSkippedSlots.Add(skippedSlots)
		avgScore += vInfo.TotalScore
	}

	dmodel.AVGScore = avgScore
	dmodel.AVGSkippedSlots = avgSkippedSlots
	if len(validators) > 0 {
		avgSkippedSlots = avgSkippedSlots.Div(decimal.New(int64(len(validators)), 0))
		avgScore = avgScore / int64(len(validators))
	}
	/*	err = s.dao.DeleteValidators(dPool.ID)
		if err != nil {
			return fmt.Errorf("dao.DeleteValidators: %s", err.Error())
		}*/

	if len(validators) > 0 {
		dmodel.Delinquent = decimal.NewFromInt(delinquent).Div(decimal.NewFromInt(int64(len(validators))))
	}

	// todo
	//dPool.UnstakeLiquidity =
	err = s.dao.UpdatePoolData(dmodel)
	if err != nil {
		return fmt.Errorf("dao.UpdatePoolData: %s", err.Error())
	}
	err = s.dao.CreateValidator(validators...)
	if err != nil {
		return fmt.Errorf("dao.CreateValidators: %s", err.Error())
	}
	return nil
}

func getNodeAddress(client *client.Client, ctx context.Context, voteAddress solana.PublicKey) (solana.PublicKey, error) {
	r, err := solana_sdk.GetVoteAccounts(client.RpcClient.Call(ctx, "getVoteAccounts", map[string]interface{}{
		"votePubkey": voteAddress.String(),
	}))
	if err != nil {
		return solana.PublicKey{}, err
	}
	if len(r.Current) > 0 {
		return solana.PublicKeyFromBase58(r.Current[0].NodePubKey)
	}
	if len(r.Delinquent) > 0 {
		return solana.PublicKeyFromBase58(r.Delinquent[0].NodePubKey)
	}

	return solana.PublicKey{}, errors.New("bad vote address")
}

func getAPY(client *client.Client, ctx context.Context, key solana.PublicKey, epochInYear float64) (decimal.Decimal, uint64, error) {
	var tes rpc.GetProgramAccountsWithContextResponse
	err := rep(func() error {
		var err error
		tes, err = client.RpcClient.GetProgramAccountsWithContextAndConfig(ctx, "Stake11111111111111111111111111111111111111",
			rpc.GetProgramAccountsConfig{
				Encoding: "base64",
				Filters: []rpc.GetProgramAccountsConfigFilter{
					{
						MemCmp: &rpc.GetProgramAccountsConfigFilterMemCmp{
							Offset: 124,
							Bytes:  key.String(),
						},
					},
				},
			},
		)
		return err
	}, 10, time.Minute*1)
	if err != nil {
		return decimal.Decimal{}, 0, err
	}

	arrAddress := make([]string, len(tes.Result.Value))
	for i, v := range tes.Result.Value {
		arrAddress[i] = v.Pubkey
	}

	var amount, balance int64

	if len(arrAddress) > 500 {
		n := int(math.Ceil(float64(len(arrAddress)) / 500))
		offset := 0
		var resp []solana_sdk.GetInflationRewardResult
		for i := 0; i < n; i++ {
			if offset+500 > len(arrAddress) {
				err = rep(func() error {
					resp, err = solana_sdk.GetInflationReward(client.RpcClient.Call(ctx, "getInflationReward", arrAddress[offset:]))
					return err
				}, 10, time.Minute*1)
				if err != nil {
					return decimal.Decimal{}, 0, err
				}
			} else {
				err = rep(func() error {
					resp, err = solana_sdk.GetInflationReward(client.RpcClient.Call(ctx, "getInflationReward", arrAddress[offset:offset+500]))
					return err
				}, 10, time.Minute*1)
				if err != nil {
					return decimal.Decimal{}, 0, err
				}
			}

			for _, v := range resp {
				amount += v.Amount
				balance += v.PostBalance
			}

			offset += 500
		}
	} else {
		var resp []solana_sdk.GetInflationRewardResult
		err = rep(func() error {
			resp, err = solana_sdk.GetInflationReward(client.RpcClient.Call(ctx, "getInflationReward", arrAddress))
			return err
		}, 10, time.Minute*1)
		if err != nil {
			return decimal.Decimal{}, 0, err
		}

		for _, v := range resp {
			amount += v.Amount
			balance += v.PostBalance
		}
	}

	if amount == 0 || balance == 0 {
		return decimal.Decimal{}, 0, nil
	}

	coefficient := decimal.NewFromInt(amount).Div(decimal.NewFromInt(balance - amount))

	return coefficient.Add(decimal.NewFromInt(1)).Pow(decimal.NewFromFloat(epochInYear)).Sub(decimal.NewFromInt(1)), uint64(len(arrAddress)), nil
}

func rep(f func() error, t uint64, timeout time.Duration) error {
	var err error
	for i := uint64(0); i < t; i++ {
		err = f()
		if err == nil {
			return nil
		}
		if i+1 < t {
			<-time.After(timeout)
		}
	}
	return err
}
