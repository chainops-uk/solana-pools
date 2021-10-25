package types

import "github.com/dfuse-io/solana-go"

var EmptyAddress = solana.MustPublicKeyFromBase58("11111111111111111111111111111111")

type (
	Pool struct {
		Address     solana.PublicKey
		SolanaStake uint64
		TokenSupply uint64
		Validators  []PoolValidator
	}
	PoolValidator struct {
		ActiveStake uint64
		NodePK      solana.PublicKey
		VotePK      solana.PublicKey
	}
)
