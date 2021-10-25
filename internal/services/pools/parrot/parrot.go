package parrot

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/dfuse-io/solana-go"
	"github.com/dfuse-io/solana-go/rpc"
	"github.com/everstake/solana-pools/internal/services/pools/types"
	"github.com/near/borsh-go"
)

type (
	Pool struct {
		solanaRPC *rpc.Client
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
		Header byte
	}
	OptionFee struct {
		Fee    Fee
		Header byte
	}
	Fee struct {
		Denominator uint64
		Numerator   uint64
	}
	Validator struct {
		ActiveStake        uint64
		TransientStake     uint64
		LastUpdateEpoch    uint64
		Status             uint8
		VoteAccountAddress solana.PublicKey
	}
)

func New(sRPC *rpc.Client) *Pool {
	return &Pool{solanaRPC: sRPC}
}

func (p Pool) GetData(address string) (data types.Pool, err error) {
	scAddress, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return data, fmt.Errorf("solana.PublicKeyFromBase58: %s", err.Error())
	}
	poolInfo, err := p.solanaRPC.GetAccountInfo(context.Background(), scAddress)
	if err != nil {
		return data, fmt.Errorf("solanaRPC.GetAccountInfo: %s", err.Error())
	}
	var poolData PoolData
	err = borsh.Deserialize(&poolData, poolInfo.Value.Data)
	if err != nil {
		return data, fmt.Errorf("borsh.Deserialize(PoolData): %s", err.Error())
	}
	valAccountInfo, err := p.solanaRPC.GetAccountInfo(context.Background(), poolData.ValidatorList)
	if err != nil {
		return data, fmt.Errorf("solanaRPC.GetAccountInfo: %s", err.Error())
	}
	var totalActiveStake uint64
	var validators []types.PoolValidator
	if len(valAccountInfo.Value.Data) > 8 { // todo
		maxValidators := binary.LittleEndian.Uint32(valAccountInfo.Value.Data[1:5])
		validatorsData := make([]Validator, maxValidators, maxValidators)
		err = borsh.Deserialize(&validatorsData, valAccountInfo.Value.Data[9:])
		if err != nil {
			return data, fmt.Errorf("borsh.Deserialize(ValidatorData): %s", err.Error())
		}
		for _, v := range validatorsData {
			if v.VoteAccountAddress.Equals(types.EmptyAddress) {
				continue
			}
			validators = append(validators, types.PoolValidator{
				ActiveStake: v.ActiveStake,
				VotePK:      v.VoteAccountAddress,
			})
			totalActiveStake += v.ActiveStake
		}
	}
	return types.Pool{
		Address:     solana.MustPublicKeyFromBase58(address),
		SolanaStake: totalActiveStake,
		TokenSupply: 0,
		Validators:  validators,
	}, nil
}
