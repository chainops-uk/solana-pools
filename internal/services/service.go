package services

import (
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/everstake/solana-pools/pkg/validatorsapp"
	"github.com/portto/solana-go-sdk/client"
	"go.uber.org/zap"
)

type (
	Service interface {
		UpdatePools() error
		GetPool(name string) (pool smodels.PoolDetails, err error)
	}
	Imp struct {
		rpcClients    map[config.Network]*client.Client
		cfg           config.Env
		dao           dao.DAO
		log           *zap.Logger
		validatorsApp *validatorsapp.Client
	}
)

func NewService(cfg config.Env, d dao.DAO, l *zap.Logger) Service {
	return &Imp{
		rpcClients: map[config.Network]*client.Client{
			config.Mainnet: client.NewClient(cfg.MainnetNode),
			config.Testnet: client.NewClient(cfg.TestnetNode),
		},
		cfg:           cfg,
		dao:           d,
		log:           l,
		validatorsApp: validatorsapp.NewClient(cfg.ValidatorsAppKey),
	}
}
