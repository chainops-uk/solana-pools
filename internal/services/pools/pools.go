package pools

import (
	"github.com/dfuse-io/solana-go/rpc"
	"github.com/everstake/solana-pools/internal/services/pools/types"
)

type (
	StakePools struct {
		solanaRPC *rpc.Client
	}
	Pool interface {
		GetData(address string) (p types.Pool, err error)
	}
)
