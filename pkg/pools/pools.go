package pools

import (
	"fmt"
	"github.com/dfuse-io/solana-go/rpc"
	"github.com/everstake/solana-pools/pkg/pools/parrot"
	"github.com/everstake/solana-pools/pkg/pools/types"
)

type (
	Pool interface {
		GetData(address string) (p types.Pool, err error)
	}
	Factory struct {
		solanaRPC *rpc.Client
	}
)

func NewFactory(rpcClient *rpc.Client) Factory {
	return Factory{solanaRPC: rpcClient}
}

func (f Factory) GetPool(name string) (p Pool, err error) {
	switch name {
	case types.ParrotPool:
		return parrot.New(f.solanaRPC), nil
	default:
		return nil, fmt.Errorf("pool %s not found", name)
	}
}
