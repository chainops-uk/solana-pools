package stdpool

import (
	"context"
	"encoding/binary"
	"fmt"
	bin "github.com/dfuse-io/binary"
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/pkg/pools/types"
	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/client"
)

type Pool struct {
	solanaRPC *client.Client
}

func New(solanaRPC *client.Client) *Pool {
	return &Pool{
		solanaRPC: solanaRPC,
	}
}

type ValidatorsData struct {
	Type       byte
	MaxSize    uint32
	Validators []ValidatorStakeInfo
}

type PoolData struct {
	AccountType                           uint8
	Manager                               solana.PublicKey
	Staker                                solana.PublicKey
	StakeDepositAuthority                 solana.PublicKey
	StakeWithdrawBumpSeed                 uint8
	ValidatorList                         solana.PublicKey
	ReserveStake                          solana.PublicKey
	PoolMint                              solana.PublicKey
	ManagerFeeAccount                     solana.PublicKey
	TokenProgramId                        solana.PublicKey
	TotalLamports                         uint64
	PoolTokenSupply                       uint64
	LastUpdateEpoch                       uint64
	Lockup                                Lockup
	EpochFee                              Fee
	NextEpochFee                          Fee
	PreferredDepositValidatorVoteAddress  solana.PublicKey
	PreferredWithdrawValidatorVoteAddress solana.PublicKey
	StakeDepositFee                       Fee
	StakeWithdrawalFee                    Fee
	NextStakeWithdrawalFee                Fee
	StakeReferralFee                      uint8
	SolDepositAuthority                   solana.PublicKey
	SolDepositFee                         Fee
	SolReferralFee                        uint8
	SolWithdrawAuthority                  solana.PublicKey
	SolWithdrawalFee                      Fee
	NextSolWithdrawalFee                  Fee
	LastEpochPoolTokenSupply              uint64
	LastEpochTotalLamports                uint64
}

type Fee struct {
	Denominator uint64
	Numerator   uint64
}

type Lockup struct {
	UnixTimestamp int64
	Epoch         uint64
	Custodian     solana.PublicKey
}

type ValidatorStakeInfo struct {
	ActiveStakeLamports      uint64
	TransientStakeLamports   int64
	LastUpdateEpoch          int64
	TransientSeedSuffixStart int64
	TransientSeedSuffixEnd   int64
	Status                   uint8
	VoteAccountAddress       solana.PublicKey
}

func validatorListFromBytes(b []byte) ([]ValidatorStakeInfo, error) {
	if len(b) > 8 { // todo
		data := ValidatorsData{}
		err := borsh.Deserialize(&data, b)
		if err != nil {
			return nil, fmt.Errorf("borsh.Deserialize(ValidatorData): %s", err.Error())
		}

		return data.Validators, err
	}
	return nil, fmt.Errorf("bad data")
}

func (sp *PoolData) SetFromBytes(b []byte) (*PoolData, error) {
	data := make([]byte, binary.Size(*sp))
	copy(data, b)
	if err := bin.NewDecoder(data).Decode(sp); err != nil {
		return nil, err
	}

	return sp, nil
}

func (p Pool) GetData(address string) (*types.Pool, error) {
	scAddress, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return nil, fmt.Errorf("solana.PublicKeyFromBase58: %s", err.Error())
	}
	poolInfo, err := p.solanaRPC.GetAccountInfo(context.Background(), scAddress.String())
	if err != nil {
		return nil, fmt.Errorf("solanaRPC.GetAccountInfo: %s", err.Error())
	}
	poolData := &PoolData{}
	poolData, err = poolData.SetFromBytes(poolInfo.Data)
	if err != nil {
		return nil, fmt.Errorf("borsh.Deserialize(PoolData): %s", err.Error())
	}
	valAccountInfo, err := p.solanaRPC.GetAccountInfo(context.Background(), poolData.ValidatorList.String())
	if err != nil {
		return nil, fmt.Errorf("solanaRPC.GetAccountInfo: %s", err.Error())
	}
	validatorsData, err := validatorListFromBytes(valAccountInfo.Data)
	if err != nil {
		return nil, fmt.Errorf("borsh.Deserialize(ValidatorData): %s", err.Error())
	}
	var totalActiveStake uint64
	validators := make([]types.PoolValidator, len(validatorsData))
	for i, v := range validatorsData {
		validators[i].ActiveStake = v.ActiveStakeLamports
		validators[i].VotePK = v.VoteAccountAddress
		totalActiveStake += v.ActiveStakeLamports
	}

	return &types.Pool{
		Address:     solana.MustPublicKeyFromBase58(address),
		SolanaStake: totalActiveStake,
		TokenSupply: poolData.PoolTokenSupply,
		Validators:  validators,
	}, nil
}
