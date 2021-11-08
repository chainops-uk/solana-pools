package solido

import (
	"context"
	"fmt"
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/pkg/pools/types"
	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/client"
)

type (
	Pool struct {
		solanaRPC *client.Client
	}
	ExchangeRate struct {
		ComputedInEpoch uint64
		StSolSupply     uint64
		SolBalance      uint64
	}
	RewardDistribution struct {
		TreasuryFee       uint32
		ValidationFee     uint32
		DeveloperFee      uint32
		StSolAppreciation uint32
	}
	FeeRecipients struct {
		TreasuryAccount  solana.PublicKey
		DeveloperAccount solana.PublicKey
	}
	Validator struct {
		PubKey solana.PublicKey

		FeeCredit              uint64
		FeeAddress             solana.PublicKey
		StakeSeeds             [16]byte
		UnstakeSeeds           [16]byte
		StakeAccountsBalance   uint64
		UnstakeAccountsBalance uint64
		Active                 bool
	}
	PoolData struct {
		Version                          byte
		Manager                          solana.PublicKey
		SolMint                          solana.PublicKey
		ExchangeRate                     ExchangeRate
		SolReserveAccountBumpSeed        byte
		StakeAuthorityBumpSeed           byte
		MintAuthorityBumpSeed            byte
		RewardsWithdrawAuthorityBumpSeed byte
		RewardDistribution               RewardDistribution
		FeeRecipients                    FeeRecipients
		Metrics                          [184]byte
		Validators                       []Validator
	}
)

func New(sRPC *client.Client) *Pool {
	return &Pool{solanaRPC: sRPC}
}

func (p Pool) GetData(address string) (data *types.Pool, err error) {
	scAddress, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return data, fmt.Errorf("solana.PublicKeyFromBase58: %s", err.Error())
	}
	poolInfo, err := p.solanaRPC.GetAccountInfo(context.Background(), scAddress.String())
	if err != nil {
		return data, fmt.Errorf("solanaRPC.GetAccountInfo: %s", err.Error())
	}
	var poolData PoolData
	err = borsh.Deserialize(&poolData, poolInfo.Data)
	if err != nil {
		return data, fmt.Errorf("borsh.Deserialize(PoolData): %s", err.Error())
	}
	var validators []types.PoolValidator
	for _, v := range poolData.Validators {
		validators = append(validators, types.PoolValidator{
			ActiveStake: v.StakeAccountsBalance,
			NodePK:      v.PubKey,
		})
	}
	rewardsFee := poolData.RewardDistribution.DeveloperFee + poolData.RewardDistribution.TreasuryFee + poolData.RewardDistribution.ValidationFee
	return &types.Pool{
		Address:       solana.MustPublicKeyFromBase58(address),
		TokenSupply:   poolData.ExchangeRate.StSolSupply,
		SolanaStake:   poolData.ExchangeRate.SolBalance,
		DepositFee:    0,
		WithdrawalFee: 0,
		RewardsFee:    float64(rewardsFee) / 100,
		Validators:    validators,
	}, nil
}
