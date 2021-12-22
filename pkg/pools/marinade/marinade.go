package marinade

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/pkg/pools/types"
	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/client"
)

/*
	https://github.com/marinade-finance/marinade-ts-cli/tree/main/src
 	+ https://project-serum.github.io/anchor/cli/commands.html#deploy
*/

type (
	Pool struct {
		solanaRPC *client.Client
	}
	StakeSystem struct {
		StakeList struct {
			Account     solana.PublicKey
			ItemSize    uint32
			Count       uint32
			NewAcc      solana.PublicKey
			CopiedCount uint32
		}
		DelayedUnstakeCoolingDown uint64
		StakeDepositBumpSeed      uint8
		StakeWithdrawBumpSeed     uint8
		SlotsForStakeDelta        uint64
		LastStakeDeltaEpoch       uint64
		MinStake                  uint64
		ExtraStakeDeltaRuns       uint32
	}

	ValidatorSystem struct {
		ValidatorList struct {
			Account     solana.PublicKey
			ItemSize    uint32
			Count       uint32
			NewAcc      solana.PublicKey
			CopiedCount uint32
		}
		ManagerAuthority        solana.PublicKey
		TotalValidatorScore     uint32
		TotalActiveBalance      uint64
		AutoAddValidatorEnabled bool
	}
	PoolData struct {
		SomeHeader                [8]byte
		MsolMint                  solana.PublicKey
		AdminAuthority            solana.PublicKey
		OperationalSolAccount     solana.PublicKey
		TreasuryMsolAccount       solana.PublicKey
		ReserveBumpSeed           byte
		MsolMintAuthorityBumpSeed byte
		RentExemptForTokenAcc     uint64
		RewardFee                 uint32
		StakeSystem               StakeSystem
		ValidatorSystem           ValidatorSystem
		LpMint                    solana.PublicKey
		LpMintAuthorityBumpSeed   byte
		SolLegBumpSeed            byte
		MsolLegAuthorityBumpSeed  byte
		MsolLeg                   solana.PublicKey
		LpLiquidityTarget         uint64
		LpMaxFee                  uint32
		LpMinFee                  uint32
		TreasuryCut               uint32
		LpSupply                  uint64
		LentFromSolLeg            uint64
		LiquiditySolCap           uint64
		AvailableReserveBalance   uint64
		MsolSupply                uint64
		MsolPrice                 uint64
		CirculatingTicketCount    uint64
		CirculatingTicketBalance  uint64
		LentFromReserve           uint64
		MinDeposit                uint64
		MinWithdraw               uint64
		StakingSolCap             uint64
		EmergencyCoolingDown      uint64
	}
	Validator struct {
		Address solana.PublicKey
		Stake   uint64
		Skip    [21]byte
	}
)

func New(client *client.Client) *Pool {
	return &Pool{
		solanaRPC: client,
	}
}

func (p Pool) GetData(address string) (*types.Pool, error) {
	scAddress, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return nil, fmt.Errorf("solana.PublicKeyFromBase58: %w", err)
	}
	poolInfo, err := p.solanaRPC.GetAccountInfo(context.Background(), scAddress.String())
	if err != nil {
		return nil, fmt.Errorf("solanaRPC.GetAccountInfo: %w", err)
	}
	var poolData PoolData
	err = borsh.Deserialize(&poolData, poolInfo.Data)
	if err != nil {
		return nil, fmt.Errorf("borsh.Deserialize(PoolData): %s", err.Error())
	}
	valAccountInfo, err := p.solanaRPC.GetAccountInfo(context.Background(), poolData.ValidatorSystem.ValidatorList.Account.String())
	if err != nil {
		return nil, fmt.Errorf("solanaRPC.GetAccountInfo: %w", err)
	}
	validatorsData := make([]*Validator, 0)

	validatorInfoSize := binary.Size(Validator{})
	for i := 8; i+validatorInfoSize <= len(valAccountInfo.Data); i += validatorInfoSize {
		d := valAccountInfo.Data[i : i+validatorInfoSize]
		var v Validator
		if err := borsh.Deserialize(&v, d); err != nil {
			return nil, fmt.Errorf("borsh.Deserialize: %w", err)
		}
		if v.Address.String() == "11111111111111111111111111111111" {
			continue
		}
		validatorsData = append(validatorsData, &v)
	}

	var totalActiveStake uint64
	var validators []types.PoolValidator
	for _, v := range validatorsData {
		validators = append(validators, types.PoolValidator{
			ActiveStake: v.Stake,
			VotePK:      v.Address,
		})
		totalActiveStake += v.Stake
	}

	return &types.Pool{
		Address:          solana.MustPublicKeyFromBase58(address),
		Epoch:            poolData.StakeSystem.LastStakeDeltaEpoch,
		APY:              1,
		SolanaStake:      totalActiveStake,
		TotalTokenSupply: poolData.MsolSupply,
		TotalLamports:    poolData.ValidatorSystem.TotalActiveBalance,
		UnstakeLiquidity: poolData.LiquiditySolCap,
		RewardsFee:       float64(poolData.RewardFee) / 100,
		WithdrawalFee:    0.03,
		Validators:       validators,
	}, nil
}
