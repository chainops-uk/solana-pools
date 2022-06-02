package services_test

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/everstake/solana-pools/internal/services/smodels"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gotest.tools/assert"
	"testing"
)

var poolArr = []*dmodels.Pool{
	{
		ID:      uuid.Must(uuid.FromString("1e85fd6d-3d32-4d86-b9a0-5ca2f7260af9")),
		Name:    "Pool1",
		Active:  true,
		CoinID:  uuid.Must(uuid.FromString("ab4cebc7-572c-4fb8-8d64-67ea5490cc7e")),
		Address: "poll_adr1",
		Network: "mainnet",
		Image:   "img",
		Coin:    *coinArr[0],
	},
	{
		ID:      uuid.Must(uuid.FromString("85353bae-ad58-4619-ba77-1ed0d2bc3143")),
		Name:    "Pool2",
		Active:  true,
		CoinID:  uuid.Must(uuid.FromString("32d63e6f-a0b6-4a49-93b3-86c688e74735")),
		Address: "poll_adr2",
		Network: "mainnet",
		Image:   "img2",
		Coin:    *coinArr[1],
	},
	{
		ID:      uuid.Must(uuid.FromString("12112f44-9a6a-4d7d-b09f-8c072a647fb2")),
		Name:    "Pool3",
		Active:  true,
		CoinID:  uuid.Must(uuid.FromString("b55a53a3-02c2-4963-81bc-385854249c96")),
		Address: "poll_adr3",
		Network: "testnet",
		Image:   "img3",
		Coin:    *coinArr[2],
	},
}

var coinArr = []*dmodels.Coin{
	{
		ID:         uuid.Must(uuid.FromString("ab4cebc7-572c-4fb8-8d64-67ea5490cc7e")),
		Name:       "coin1",
		GeckoKey:   "key1",
		Address:    "coin_addr1",
		USD:        73,
		ThumbImage: "img1",
		SmallImage: "none",
		LargeImage: "none",
	},
	{
		ID:         uuid.Must(uuid.FromString("32d63e6f-a0b6-4a49-93b3-86c688e74735")),
		Name:       "coin2",
		GeckoKey:   "key2",
		Address:    "coin_addr2",
		USD:        105,
		ThumbImage: "img2",
		SmallImage: "none",
		LargeImage: "none",
	},
	{
		ID:         uuid.Must(uuid.FromString("644740ed-3b5c-4a6c-bd3b-45003ce852a1")),
		Name:       "coin3",
		GeckoKey:   "key3",
		Address:    "coin_addr3",
		USD:        81,
		ThumbImage: "img3",
		SmallImage: "none",
		LargeImage: "none",
	},
}

var DeFiArr = []*dmodels.DEFI{
	{
		ID:              uuid.Must(uuid.FromString("89dc136c-7d8e-46eb-9329-f3823435ba01")),
		LiquidityPoolID: uuid.Must(uuid.FromString("721dd49b-0e19-4655-9052-42c8a57aef01")),
		SaleCoinID:      uuid.Must(uuid.FromString("fa3f2eb5-eccd-47b9-b837-bdaa8d825ac6")),
		BuyCoinID:       uuid.Must(uuid.FromString("9002fd64-4e32-4555-91e7-1a2399faf9cc")),
		Liquidity:       53,
		APY:             decimal.NewFromInt32(50),
		SaleCoin:        dmodels.Coin{},
		BuyCoin:         dmodels.Coin{},
		LiquidityPool:   LPArr[0],
	},
}

var LPArr = []dmodels.LiquidityPool{
	{
		ID:    uuid.Must(uuid.FromString("721dd49b-0e19-4655-9052-42c8a57aef01")),
		Name:  "LPName1",
		About: "123fdg",
		Image: "LPImg1",
		URL:   "LPUrl1",
	},
	{
		ID:    uuid.Must(uuid.FromString("0a176030-4624-4da9-9d13-6991ec379d1c")),
		Name:  "LPName2",
		About: "dfsdfg",
		Image: "LPImg2",
		URL:   "LPUrl2",
	},
	{
		ID:    uuid.Must(uuid.FromString("756e4171-14c9-494a-af4a-e01b5ba69e5f")),
		Name:  "LPName3",
		About: "34ffsd",
		Image: "LPImg3",
		URL:   "LPUrl3",
	},
}

func TestGetPoolCoins(t *testing.T) {
	data := map[string]struct {
		DAO  services.Imp
		Data struct {
			name   string
			sort   string
			desc   bool
			limit  uint64
			offset uint64
		}
		Result []*smodels.Coin
		Err    error
	}{
		"first": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "coin1", sort: "price", desc: true, limit: 2, offset: 0},
			Result: []*smodels.Coin{
				{
					Name:       "coin1",
					Address:    "coin_addr1",
					USD:        73,
					ThumbImage: "img1",
					SmallImage: "none",
					LargeImage: "none",
					DeFi: []*smodels.DeFi{
						(&smodels.DeFi{}).Set(DeFiArr[0], (&smodels.Coin{}).Set(coinArr[0], nil), (&smodels.LiquidityPool{}).Set(&LPArr[0])),
					},
				},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.name != %s, name is %s", postgres.MainNet, condition.Network)
						}

						return poolArr[:2], nil
					},
					GetCoinsFunc: func(cond *postgres.CoinCondition) ([]*dmodels.Coin, error) {
						if cond.Condition.IDs[0] != poolArr[0].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[0] != %s, id is %s", poolArr[0].Coin.ID, cond.Condition.IDs[0])
						}
						if cond.Condition.IDs[1] != poolArr[1].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[1] != %s, id is %s", poolArr[1].Coin.ID, cond.Condition.IDs[1])
						}
						if cond.Condition.Pagination.Limit != 2 {
							return nil, fmt.Errorf("limit != 500, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Condition.Pagination.Offset)
						}
						if cond.CoinSort.Sort != postgres.CoinPrice {
							return nil, fmt.Errorf("cond.CoinSort.Sort != postgres.CoinPrice, sort = %d", cond.CoinSort.Sort)
						}
						if cond.CoinSort.Desc != true {
							return nil, fmt.Errorf("cond.CoinSort.Desc != true, desc = %v", false)
						}
						if cond.Name != "coin1" {
							return nil, fmt.Errorf("cond.Name != coin1, name = %s", cond.Name)
						}
						return coinArr[:1], nil
					},
					GetDEFIsFunc: func(cond *postgres.DeFiCondition) ([]*dmodels.DEFI, error) {
						if cond.SaleCoinID[0] != coinArr[0].ID {
							return nil, fmt.Errorf("cond.SaleCoinID[0] != coinArr[0].ID, cond.SaleCoinID[0] = %s", cond.SaleCoinID[0])
						}
						return DeFiArr[:1], nil
					},
					GetLiquidityPoolFunc: func(cond *postgres.Condition) (*dmodels.LiquidityPool, error) {
						if cond.IDs[0] != LPArr[0].ID {
							return nil, fmt.Errorf("cond.IDs[0] != LPArr[0].ID, cond.IDs[0] = %s", cond.IDs[0])
						}
						return &LPArr[0], nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != DeFiArr[0].BuyCoinID {
							return nil, fmt.Errorf("id != DeFiArr[0].BuyCoinID, id = %v", id)
						}
						return coinArr[0], nil
					},
					GetCoinsCountFunc: func(cond *postgres.CoinCondition) (int64, error) {
						if cond.Condition.IDs[0] != poolArr[0].Coin.ID {
							return 0, fmt.Errorf("cond.Condition.IDs[0] != %s, cond.Condition.IDs[0] = %s", poolArr[0].Coin.ID, cond.Condition.IDs[0])
						}
						if cond.Condition.IDs[1] != poolArr[1].Coin.ID {
							return 0, fmt.Errorf("cond.Condition.IDs[1] != %s, cond.Condition.IDs[0] = %s", poolArr[1].Coin.ID, cond.Condition.IDs[1])
						}
						if cond.Name != "coin1" {
							return 0, fmt.Errorf("cond.Name != coin1, cond.Name = %s", cond.Name)
						}
						return 1, nil
					},
				},
			},
		},
		"second": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "coin5", sort: "price", desc: true, limit: 500, offset: 1},
			Result: []*smodels.Coin{},
			Err:    nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.name != %s, name is %s", postgres.MainNet, condition.Network)
						}

						return poolArr[:2], nil
					},
					GetCoinsFunc: func(cond *postgres.CoinCondition) ([]*dmodels.Coin, error) {
						if cond.Condition.IDs[0] != poolArr[0].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[0] != %s, id is %s", poolArr[0].Coin.ID, cond.Condition.IDs[0])
						}
						if cond.Condition.IDs[1] != poolArr[1].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[1] != %s, id is %s", poolArr[0].Coin.ID, cond.Condition.IDs[1])
						}
						if cond.Condition.Pagination.Limit != 500 {
							return nil, fmt.Errorf("limit != 500, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 1 {
							return nil, fmt.Errorf("offset != 500, offset = %d", cond.Condition.Pagination.Offset)
						}
						if cond.CoinSort.Sort != postgres.CoinPrice {
							return nil, fmt.Errorf("cond.CoinSort.Sort != postgres.CoinPrice, sort = %d", cond.CoinSort.Sort)
						}
						if cond.CoinSort.Desc != true {
							return nil, fmt.Errorf("cond.CoinSort.Desc != true, desc = %v", false)
						}
						if cond.Name != "coin5" {
							return nil, fmt.Errorf("cond.Name != coin5, name = %s", cond.Name)
						}
						return nil, nil
					},
					GetDEFIsFunc: func(cond *postgres.DeFiCondition) ([]*dmodels.DEFI, error) {
						return nil, nil
					},
					GetLiquidityPoolFunc: func(cond *postgres.Condition) (*dmodels.LiquidityPool, error) {
						return nil, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						return nil, nil
					},
					GetCoinsCountFunc: func(cond *postgres.CoinCondition) (int64, error) {
						if cond.Condition.IDs[0] != poolArr[0].Coin.ID {
							return 0, fmt.Errorf("cond.Condition.IDs[0] != %s, cond.Condition.IDs[0] = %s", poolArr[0].Coin.ID, cond.Condition.IDs[0])
						}
						if cond.Condition.IDs[1] != poolArr[1].Coin.ID {
							return 0, fmt.Errorf("cond.Condition.IDs[1] != %s, cond.Condition.IDs[0] = %s", poolArr[1].Coin.ID, cond.Condition.IDs[1])
						}
						if cond.Name != "coin5" {
							return 0, fmt.Errorf("cond.Name != coin5, cond.Name = %s", cond.Name)
						}
						return 0, nil
					},
				},
			},
		},
		"third": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "coin3", sort: "price", desc: true, limit: 500, offset: 1},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetCoins: %w", gorm.ErrRecordNotFound),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.name != %s, name is %s", postgres.MainNet, condition.Network)
						}

						return poolArr[:2], nil
					},
					GetCoinsFunc: func(cond *postgres.CoinCondition) ([]*dmodels.Coin, error) {
						if cond.Condition.IDs[0] != poolArr[0].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[0] != %s, id is %s", poolArr[0].Coin.ID, cond.Condition.IDs[0])
						}
						if cond.Condition.IDs[1] != poolArr[1].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[1] != %s, id is %s", poolArr[1].Coin.ID, cond.Condition.IDs[1])
						}
						if cond.Condition.Pagination.Limit != 500 {
							return nil, fmt.Errorf("limit != 500, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 1 {
							return nil, fmt.Errorf("offset != 500, offset = %d", cond.Condition.Pagination.Offset)
						}
						if cond.CoinSort.Sort != postgres.CoinPrice {
							return nil, fmt.Errorf("cond.CoinSort.Sort != postgres.CoinPrice, sort = %d", cond.CoinSort.Sort)
						}
						if cond.CoinSort.Desc != true {
							return nil, fmt.Errorf("cond.CoinSort.Desc != true, desc = %v", false)
						}
						if cond.Name != "coin3" {
							return nil, fmt.Errorf("cond.Name != coin3, name = %s", cond.Name)
						}
						return nil, gorm.ErrRecordNotFound
					},
					GetDEFIsFunc: func(cond *postgres.DeFiCondition) ([]*dmodels.DEFI, error) {
						if cond.SaleCoinID[0] != coinArr[0].ID {
							return nil, fmt.Errorf("cond.SaleCoinID[0] != coinArr[0].ID, cond.SaleCoinID[0] = %s", cond.SaleCoinID[0])
						}
						return DeFiArr[:1], nil
					},
					GetLiquidityPoolFunc: func(cond *postgres.Condition) (*dmodels.LiquidityPool, error) {
						if cond.IDs[0] != LPArr[0].ID {
							return nil, fmt.Errorf("cond.IDs[0] != LPArr[0].ID, cond.IDs[0] = %s", cond.IDs[0])
						}
						return &LPArr[0], nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != DeFiArr[0].BuyCoinID {
							return nil, fmt.Errorf("id != DeFiArr[0].BuyCoinID, id = %v", id)
						}
						return coinArr[0], nil
					},
					GetCoinsCountFunc: func(cond *postgres.CoinCondition) (int64, error) {
						if cond.Condition.IDs[0] != poolArr[0].Coin.ID {
							return 0, fmt.Errorf("cond.Condition.IDs[0] != %s, cond.Condition.IDs[0] = %s", poolArr[0].Coin.ID, cond.Condition.IDs[0])
						}
						if cond.Condition.IDs[1] != poolArr[1].Coin.ID {
							return 0, fmt.Errorf("cond.Condition.IDs[1] != %s, cond.Condition.IDs[0] = %s", poolArr[1].Coin.ID, cond.Condition.IDs[1])
						}
						if cond.Name != "coin1" {
							return 0, fmt.Errorf("cond.Name != coin1, cond.Name = %s", cond.Name)
						}
						return 1, nil
					},
				},
			},
		},
		"forth": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "coin1", sort: "price", desc: true, limit: 2, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetDEFIs: %w", gorm.ErrRecordNotFound),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.name != %s, name is %s", postgres.MainNet, condition.Network)
						}

						return poolArr[:2], nil
					},
					GetCoinsFunc: func(cond *postgres.CoinCondition) ([]*dmodels.Coin, error) {
						if cond.Condition.IDs[0] != poolArr[0].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[0] != %s, id is %s", poolArr[0].Coin.ID, cond.Condition.IDs[0])
						}
						if cond.Condition.IDs[1] != poolArr[1].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[1] != %s, id is %s", poolArr[1].Coin.ID, cond.Condition.IDs[1])
						}
						if cond.Condition.Pagination.Limit != 2 {
							return nil, fmt.Errorf("limit != 500, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Condition.Pagination.Offset)
						}
						if cond.CoinSort.Sort != postgres.CoinPrice {
							return nil, fmt.Errorf("cond.CoinSort.Sort != postgres.CoinPrice, sort = %d", cond.CoinSort.Sort)
						}
						if cond.CoinSort.Desc != true {
							return nil, fmt.Errorf("cond.CoinSort.Desc != true, desc = %v", false)
						}
						if cond.Name != "coin1" {
							return nil, fmt.Errorf("cond.Name != coin1, name = %s", cond.Name)
						}
						return coinArr[:1], nil
					},
					GetDEFIsFunc: func(cond *postgres.DeFiCondition) ([]*dmodels.DEFI, error) {
						return nil, gorm.ErrRecordNotFound
					},
				},
			},
		},
		"fifth": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "coin1", sort: "price", desc: true, limit: 2, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetLiquidityPool: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.name != %s, name is %s", postgres.MainNet, condition.Network)
						}

						return poolArr[:2], nil
					},
					GetCoinsFunc: func(cond *postgres.CoinCondition) ([]*dmodels.Coin, error) {
						if cond.Condition.IDs[0] != poolArr[0].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[0] != %s, id is %s", poolArr[0].Coin.ID, cond.Condition.IDs[0])
						}
						if cond.Condition.IDs[1] != poolArr[1].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[1] != %s, id is %s", poolArr[1].Coin.ID, cond.Condition.IDs[1])
						}
						if cond.Condition.Pagination.Limit != 2 {
							return nil, fmt.Errorf("limit != 500, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Condition.Pagination.Offset)
						}
						if cond.CoinSort.Sort != postgres.CoinPrice {
							return nil, fmt.Errorf("cond.CoinSort.Sort != postgres.CoinPrice, sort = %d", cond.CoinSort.Sort)
						}
						if cond.CoinSort.Desc != true {
							return nil, fmt.Errorf("cond.CoinSort.Desc != true, desc = %v", false)
						}
						if cond.Name != "coin1" {
							return nil, fmt.Errorf("cond.Name != coin1, name = %s", cond.Name)
						}
						return coinArr[:1], nil
					},
					GetDEFIsFunc: func(cond *postgres.DeFiCondition) ([]*dmodels.DEFI, error) {
						if cond.SaleCoinID[0] != coinArr[0].ID {
							return nil, fmt.Errorf("cond.SaleCoinID[0] != coinArr[0].ID, cond.SaleCoinID[0] = %s", cond.SaleCoinID[0])
						}
						return DeFiArr[:1], nil
					},
					GetLiquidityPoolFunc: func(cond *postgres.Condition) (*dmodels.LiquidityPool, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"sixth": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "coin1", sort: "price", desc: true, limit: 2, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetCoinByID: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.name != %s, name is %s", postgres.MainNet, condition.Network)
						}

						return poolArr[:2], nil
					},
					GetCoinsFunc: func(cond *postgres.CoinCondition) ([]*dmodels.Coin, error) {
						if cond.Condition.IDs[0] != poolArr[0].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[0] != %s, id is %s", poolArr[0].Coin.ID, cond.Condition.IDs[0])
						}
						if cond.Condition.IDs[1] != poolArr[1].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[1] != %s, id is %s", poolArr[1].Coin.ID, cond.Condition.IDs[1])
						}
						if cond.Condition.Pagination.Limit != 2 {
							return nil, fmt.Errorf("limit != 500, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Condition.Pagination.Offset)
						}
						if cond.CoinSort.Sort != postgres.CoinPrice {
							return nil, fmt.Errorf("cond.CoinSort.Sort != postgres.CoinPrice, sort = %d", cond.CoinSort.Sort)
						}
						if cond.CoinSort.Desc != true {
							return nil, fmt.Errorf("cond.CoinSort.Desc != true, desc = %v", false)
						}
						if cond.Name != "coin1" {
							return nil, fmt.Errorf("cond.Name != coin1, name = %s", cond.Name)
						}
						return coinArr[:1], nil
					},
					GetDEFIsFunc: func(cond *postgres.DeFiCondition) ([]*dmodels.DEFI, error) {
						if cond.SaleCoinID[0] != coinArr[0].ID {
							return nil, fmt.Errorf("cond.SaleCoinID[0] != coinArr[0].ID, cond.SaleCoinID[0] = %s", cond.SaleCoinID[0])
						}
						return DeFiArr[:1], nil
					},
					GetLiquidityPoolFunc: func(cond *postgres.Condition) (*dmodels.LiquidityPool, error) {
						return nil, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"seventh": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "coin1", sort: "price", desc: true, limit: 2, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetCoinsCount: %w", gorm.ErrRecordNotFound),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.name != %s, name is %s", postgres.MainNet, condition.Network)
						}

						return poolArr[:2], nil
					},
					GetCoinsFunc: func(cond *postgres.CoinCondition) ([]*dmodels.Coin, error) {
						if cond.Condition.IDs[0] != poolArr[0].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[0] != %s, id is %s", poolArr[0].Coin.ID, cond.Condition.IDs[0])
						}
						if cond.Condition.IDs[1] != poolArr[1].Coin.ID {
							return nil, fmt.Errorf("condition.Condition.IDs[1] != %s, id is %s", poolArr[1].Coin.ID, cond.Condition.IDs[1])
						}
						if cond.Condition.Pagination.Limit != 2 {
							return nil, fmt.Errorf("limit != 500, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Condition.Pagination.Offset)
						}
						if cond.CoinSort.Sort != postgres.CoinPrice {
							return nil, fmt.Errorf("cond.CoinSort.Sort != postgres.CoinPrice, sort = %d", cond.CoinSort.Sort)
						}
						if cond.CoinSort.Desc != true {
							return nil, fmt.Errorf("cond.CoinSort.Desc != true, desc = %v", false)
						}
						if cond.Name != "coin1" {
							return nil, fmt.Errorf("cond.Name != coin1, name = %s", cond.Name)
						}
						return coinArr[:1], nil
					},
					GetDEFIsFunc: func(cond *postgres.DeFiCondition) ([]*dmodels.DEFI, error) {
						if cond.SaleCoinID[0] != coinArr[0].ID {
							return nil, fmt.Errorf("cond.SaleCoinID[0] != coinArr[0].ID, cond.SaleCoinID[0] = %s", cond.SaleCoinID[0])
						}
						return DeFiArr[:1], nil
					},
					GetLiquidityPoolFunc: func(cond *postgres.Condition) (*dmodels.LiquidityPool, error) {
						if cond.IDs[0] != LPArr[0].ID {
							return nil, fmt.Errorf("cond.IDs[0] != LPArr[0].ID, cond.IDs[0] = %s", cond.IDs[0])
						}
						return &LPArr[0], nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != DeFiArr[0].BuyCoinID {
							return nil, fmt.Errorf("id != DeFiArr[0].BuyCoinID, id = %v", id)
						}
						return coinArr[0], nil
					},
					GetCoinsCountFunc: func(cond *postgres.CoinCondition) (int64, error) {
						return 0, gorm.ErrRecordNotFound
					},
				},
			},
		},
	}

	for s, s2 := range data {
		t.Run(s, func(t *testing.T) {
			coins, count, err := s2.DAO.GetPoolCoins(s2.Data.name, s2.Data.sort, s2.Data.desc, s2.Data.limit, s2.Data.offset)
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			assert.Equal(t, uint64(len(s2.Result)), count)
			assert.Equal(t, uint64(len(coins)), count)
			for i, coin := range coins {
				t.Run(fmt.Sprintf("coins[%d]", i), func(t *testing.T) {
					assert.DeepEqual(t, coin, s2.Result[i])
				})
			}
		})
	}
}

func TestGetCoins(t *testing.T) {
	data := map[string]struct {
		DAO  services.Imp
		Data struct {
			name   string
			limit  uint64
			offset uint64
		}
		Result []*smodels.Coin
		Err    error
	}{
		"first": {
			Data: struct {
				name   string
				limit  uint64
				offset uint64
			}{name: "coin1", limit: 2, offset: 0},
			Result: []*smodels.Coin{
				{
					Name:       "coin1",
					Address:    "coin_addr1",
					USD:        73,
					ThumbImage: "img1",
					SmallImage: "none",
					LargeImage: "none",
					DeFi:       nil,
				},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetCoinsFunc: func(cond *postgres.CoinCondition) ([]*dmodels.Coin, error) {
						if cond.Condition.Name != "coin1" {
							return nil, fmt.Errorf("cond.Condition.Name != coin1, name = %s", cond.Name)
						}
						if cond.Condition.Pagination.Limit != 2 {
							return nil, fmt.Errorf("limit != 2, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Condition.Pagination.Offset)
						}
						return coinArr[:1], nil
					},
					GetCoinsCountFunc: func(cond *postgres.CoinCondition) (int64, error) {
						if cond.Name != "coin1" {
							return 0, fmt.Errorf("cond.Name != coin1, cond.Name = %s", cond.Name)
						}
						return 1, nil
					},
				},
			},
		},
		"second": {
			Data: struct {
				name   string
				limit  uint64
				offset uint64
			}{name: "coin1", limit: 2, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetCoins: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetCoinsFunc: func(cond *postgres.CoinCondition) ([]*dmodels.Coin, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"third": {
			Data: struct {
				name   string
				limit  uint64
				offset uint64
			}{name: "coin1", limit: 2, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetCoinsCount: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetCoinsFunc: func(cond *postgres.CoinCondition) ([]*dmodels.Coin, error) {
						if cond.Condition.Name != "coin1" {
							return nil, fmt.Errorf("cond.Condition.Name != coin1, name = %s", cond.Name)
						}
						if cond.Condition.Pagination.Limit != 2 {
							return nil, fmt.Errorf("limit != 2, limit = %d", cond.Condition.Pagination.Limit)
						}
						if cond.Condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("offset != 0, offset = %d", cond.Condition.Pagination.Offset)
						}
						return coinArr[:1], nil
					},
					GetCoinsCountFunc: func(cond *postgres.CoinCondition) (int64, error) {
						return 0, fmt.Errorf("some error")
					},
				},
			},
		},
	}

	for s, s2 := range data {
		t.Run(s, func(t *testing.T) {
			coins, count, err := s2.DAO.GetCoins(s2.Data.name, s2.Data.limit, s2.Data.offset)
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			assert.Equal(t, uint64(len(s2.Result)), count)
			assert.Equal(t, uint64(len(coins)), count)
			for i, coin := range coins {
				t.Run(fmt.Sprintf("coins[%d]", i), func(t *testing.T) {
					assert.DeepEqual(t, coin, s2.Result[i])
				})
			}
		})
	}
}
