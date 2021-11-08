package pools

import (
	"fmt"
	"github.com/everstake/solana-pools/pkg/pools/marinade"
	"github.com/everstake/solana-pools/pkg/pools/parrot"
	"github.com/everstake/solana-pools/pkg/pools/solido"
	"github.com/everstake/solana-pools/pkg/pools/stdpool"
	"github.com/everstake/solana-pools/pkg/pools/types"
	"github.com/portto/solana-go-sdk/client"
)

type (
	Pool interface {
		GetData(address string) (p *types.Pool, err error)
	}
	Factory struct {
		solanaRPC *client.Client
	}
)

func NewFactory(rpcClient *client.Client) Factory {
	return Factory{solanaRPC: rpcClient}
}

func (f Factory) GetPool(name string) (p Pool, err error) {
	switch name {
	case types.ParrotPool:
		return parrot.New(f.solanaRPC), nil
	case types.MarinadePool:
		return marinade.New(f.solanaRPC), nil
	case types.SolidoPool:
		return solido.New(f.solanaRPC), nil
	case types.EverSOL, types.Socean, types.JPool:
		return stdpool.New(f.solanaRPC), nil
	default:
		return nil, fmt.Errorf("pool %s not found", name)
	}
}
