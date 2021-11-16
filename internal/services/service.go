package services

import (
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao"
	"github.com/everstake/solana-pools/internal/dao/cache"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/everstake/solana-pools/pkg/validatorsapp"
	"github.com/portto/solana-go-sdk/client"
	"github.com/shopspring/decimal"
	coingecko "github.com/superoo7/go-gecko/v3"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type (
	PoolConditional struct {
		Name string
	}

	Service interface {
		GetPoolCount() (int64, error)
		GetActiveStake() uint64
		GetPoolsStatistic() (*smodels.Statistic, error)
		GetPrice() (decimal.Decimal, error)
		GetAPY() (decimal.Decimal, error)
		GetValidators() (int64, error)
		GetPool(name string) (pool smodels.PoolDetails, err error)
		GetPools(string, uint64, uint64) ([]*smodels.PoolDetails, error)
		UpdatePrice() error
		UpdatePools() error
		UpdateAPY() error
		UpdateValidators() error
	}
	Imp struct {
		rpcClients    map[config.Network]*client.Client
		cache         *cache.Cache
		cfg           config.Env
		dao           dao.DAO
		coinGecko     *coingecko.Client
		log           *zap.Logger
		validatorsApp *validatorsapp.Client
	}
)

func NewService(cfg config.Env, d dao.DAO, l *zap.Logger) Service {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return &Imp{
		rpcClients: map[config.Network]*client.Client{
			config.Mainnet: client.NewClient(cfg.MainnetNode),
			config.Testnet: client.NewClient(cfg.TestnetNode),
		},
		cache:         cache.New(time.Hour*24, time.Hour*24),
		cfg:           cfg,
		dao:           d,
		coinGecko:     coingecko.NewClient(httpClient),
		log:           l,
		validatorsApp: validatorsapp.NewClient(cfg.ValidatorsAppKey),
	}
}
