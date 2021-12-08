package services

import (
	"context"
	"github.com/dfuse-io/solana-go"
	solana_sdk "github.com/everstake/solana-pools/pkg/extension/solana-sdk"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/shopspring/decimal"
	"math"
	"time"
)

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
