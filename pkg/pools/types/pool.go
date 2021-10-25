package types

import "github.com/dfuse-io/solana-go"

var EmptyAddress = solana.MustPublicKeyFromBase58("11111111111111111111111111111111")

const (
	ParrotPool   = "parrot"
	MarinadePool = "marinade"
	SolidoPool   = "solido"
)

type (
	Pool struct {
		Address       solana.PublicKey
		SolanaStake   uint64
		TokenSupply   uint64
		DepositFee    float64
		WithdrawalFee float64
		RewardsFee    float64
		Validators    []PoolValidator
	}
	PoolValidator struct {
		ActiveStake uint64
		NodePK      solana.PublicKey
		VotePK      solana.PublicKey
	}
)
