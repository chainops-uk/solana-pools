package services

import (
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao"
	"github.com/everstake/solana-pools/internal/dao/cache"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/everstake/solana-pools/pkg/atrix"
	"github.com/everstake/solana-pools/pkg/orca"
	"github.com/everstake/solana-pools/pkg/raydium"
	"github.com/everstake/solana-pools/pkg/saber"
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
		GetActiveStake() uint64
		GetPoolsCurrentStatistic() (*smodels.Statistic, error)
		GetPoolStatistic(name string, aggregate string) ([]*smodels.Pool, error)
		GetPrice() (decimal.Decimal, error)
		GetAPY() (decimal.Decimal, error)
		GetValidators() (int64, error)
		GetPool(name string) (*smodels.PoolDetails, error)
		GetPools(name string, sort string, desc bool, from uint64, to uint64) ([]*smodels.PoolDetails, uint64, error)
		GetEpoch() (*smodels.EpochInfo, error)
		GetPoolCoins(name string, sort string, desc bool, limit uint64, offset uint64) ([]*smodels.Coin, uint64, error)
		GetGovernance(name string, sort string, desc bool, limit uint64, offset uint64) ([]*smodels.Governance, uint64, error)
		GetCoins(name string, limit uint64, offset uint64) ([]*smodels.Coin, uint64, error)
		GetAllValidators(validatorName string, sort string, desc bool, limit uint64, offset uint64) ([]*smodels.Validator, uint64, error)
		GetPoolValidators(name string, validatorName string, sort string, desc bool, limit uint64, offset uint64) ([]*smodels.ValidatorData, uint64, error)
		GetLiquidityPools(name string, limit uint64, offset uint64) ([]*smodels.LiquidityPool, uint64, error)
		GetAvgSlotTimeMS() (float64, error)

		UpdateDeFi() error
		UpdateCoins() error
		UpdateGovernance() error
		UpdatePrice() error
		UpdatePools() error
		UpdateNetworkData() error
		UpdateValidators() error
		UpdateSlotTimeMS() error
	}
	Imp struct {
		rpcClients    map[config.Network]*client.Client
		cache         *cache.Cache
		cfg           config.Env
		dao           dao.DAO
		coinGecko     *coingecko.Client
		log           *zap.Logger
		raydium       *raydium.Client
		atrix         *atrix.Client
		orca          *orca.Client
		saber         *saber.Client
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
		raydium:       raydium.NewClient(httpClient),
		orca:          orca.NewClient(httpClient),
		saber:         saber.NewClient(httpClient),
		atrix:         atrix.NewClient(httpClient),
		coinGecko:     coingecko.NewClient(httpClient),
		log:           l,
		validatorsApp: validatorsapp.NewClient(cfg.ValidatorsAppKey),
	}
}
