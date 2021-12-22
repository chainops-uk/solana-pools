package parrot

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
	PoolData struct {
		AccountType                           uint8
		Manager                               solana.PublicKey
		Staker                                solana.PublicKey
		StakeDepositAuthority                 solana.PublicKey
		StakeWithdrawBumpSeed                 byte
		ValidatorList                         solana.PublicKey
		ReserveStake                          solana.PublicKey
		PoolMint                              solana.PublicKey
		ManagerFee                            solana.PublicKey
		TokenProgramId                        solana.PublicKey
		TotalLamports                         uint64
		PoolTokenSupply                       uint64
		LastUpdateEpoch                       uint64
		Lockup                                solana.PublicKey
		EpochFee                              Fee
		NextEpochFee                          OptionFee
		PreferredDepositValidatorVoteAddress  OptionPubKey
		PreferredWithdrawValidatorVoteAddress OptionPubKey
		StakeDepositFee                       Fee
		StakeWithdrawalFee                    Fee
		NextStakeWithdrawalFee                OptionFee
		StakeReferralFee                      uint8
		SolDepositAuthority                   OptionPubKey
		SolDepositFee                         Fee
		SolReferralFee                        uint8
		SolWithdrawAuthority                  OptionPubKey
		SolWithdrawalFee                      Fee
		NextSolWithdrawalFee                  OptionFee
		LastEpochPoolTokenSupply              uint64
		LastEpochTotalLamports                uint64
	}
	OptionPubKey struct {
		PK     solana.PublicKey
		Footer byte
	}
	OptionFee struct {
		Fee    Fee
		Footer byte
	}
	Fee struct {
		Denominator uint64
		Numerator   uint64
	}
	ValidatorsData struct {
		Type       byte
		MaxSize    uint32
		Validators []Validator
	}
	Validator struct {
		ActiveStake        uint64
		TransientStake     uint64
		LastUpdateEpoch    uint64
		Status             uint8
		VoteAccountAddress solana.PublicKey
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
	valAccountInfo, err := p.solanaRPC.GetAccountInfo(context.Background(), poolData.ValidatorList.String())
	if err != nil {
		return data, fmt.Errorf("solanaRPC.GetAccountInfo: %s", err.Error())
	}
	var validatorsData ValidatorsData
	err = borsh.Deserialize(&validatorsData, valAccountInfo.Data)
	if err != nil {
		return data, fmt.Errorf("borsh.Deserialize(ValidatorData): %s", err.Error())
	}
	var totalActiveStake uint64
	var validators []types.PoolValidator
	for _, v := range validatorsData.Validators {
		validators = append(validators, types.PoolValidator{
			ActiveStake: v.ActiveStake,
			VotePK:      v.VoteAccountAddress,
		})
		totalActiveStake += v.ActiveStake
	}
	var depositFee, withdrawalFee, rewardsFee float64
	//if poolData.SolDepositFee.Denominator != 0 {
	//	depositFee = float64(poolData.SolDepositFee.Numerator) / float64(poolData.SolDepositFee.Denominator)
	//}
	if poolData.NextSolWithdrawalFee.Fee.Denominator != 0 {
		withdrawalFee = float64(poolData.NextSolWithdrawalFee.Fee.Numerator) / float64(poolData.NextSolWithdrawalFee.Fee.Denominator)
	}
	if poolData.EpochFee.Denominator != 0 {
		rewardsFee = float64(poolData.EpochFee.Numerator) / float64(poolData.EpochFee.Denominator)
	}

	l, err := p.solanaRPC.GetBalance(context.Background(), poolData.ReserveStake.String())
	if err != nil {
		return nil, fmt.Errorf("client.GetBalance: %s", err.Error())
	}

	_ = depositFee
	return &types.Pool{
		Address:          solana.MustPublicKeyFromBase58(address),
		Epoch:            poolData.LastUpdateEpoch,
		TotalLamports:    poolData.TotalLamports,
		SolanaStake:      totalActiveStake,
		TotalTokenSupply: poolData.PoolTokenSupply,
		UnstakeLiquidity: l,
		DepositFee:       0.1,
		WithdrawalFee:    withdrawalFee,
		RewardsFee:       rewardsFee,
		Validators:       validators,
	}, nil
}
