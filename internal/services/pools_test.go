package services_test

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao"
	"github.com/everstake/solana-pools/internal/dao/cache"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/everstake/solana-pools/pkg/models/sol"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gotest.tools/assert"
	"testing"
	"time"
)

var dPool = dmodels.Pool{
	ID:      uuid.Must(uuid.FromString("6aeb256a-55c9-450c-94b8-e3029eab0ed3")),
	Name:    "pool1",
	Active:  true,
	CoinID:  uuid.Must(uuid.FromString("ab4cebc7-572c-4fb8-8d64-67ea5490cc7e")),
	Address: "addr1",
	Network: "ntw1",
	Image:   "img1",
	Coin:    dmodels.Coin{},
}

var dPools = []dmodels.Pool{
	{
		ID:      uuid.Must(uuid.FromString("6aeb256a-55c9-450c-94b8-e3029eab0ed3")),
		Name:    "pool1",
		Active:  true,
		CoinID:  uuid.Must(uuid.FromString("ab4cebc7-572c-4fb8-8d64-67ea5490cc7e")),
		Address: "addr1",
		Network: "ntw1",
		Image:   "img1",
		Coin:    dmodels.Coin{},
	},
	{
		ID:      uuid.Must(uuid.FromString("4b70cb5a-4289-4d41-afec-b71f697cb82e")),
		Name:    "pool2",
		Active:  true,
		CoinID:  uuid.Must(uuid.FromString("ab4cebc7-572c-4fb8-8d64-67ea5490cc7e")),
		Address: "addr2",
		Network: "ntw2",
		Image:   "img2",
		Coin:    dmodels.Coin{},
	},
	{
		ID:      uuid.Must(uuid.FromString("c0dce5ca-8f32-4fe5-9c33-dae5dd4940ea")),
		Name:    "pool3",
		Active:  true,
		CoinID:  uuid.Must(uuid.FromString("ab4cebc7-572c-4fb8-8d64-67ea5490cc7e")),
		Address: "addr3",
		Network: "ntw3",
		Image:   "img3",
		Coin:    dmodels.Coin{},
	},
}

var poolDetails = smodels.PoolDetails{
	Pool: smodels.Pool{
		Address:          "addr1",
		Name:             "pool1",
		Image:            "img1",
		Currency:         "cur1",
		ActiveStake:      sol.SOL{},
		TokensSupply:     sol.SOL{},
		TotalLamports:    sol.SOL{},
		APY:              decimal.Decimal{},
		AVGSkippedSlots:  decimal.Decimal{},
		AVGScore:         10,
		StakingAccounts:  20,
		Delinquent:       30,
		UnstakeLiquidity: sol.SOL{},
		DepossitFee:      decimal.Decimal{},
		WithdrawalFee:    decimal.Decimal{},
		RewardsFee:       decimal.Decimal{},
		ValidatorCount:   40,
		CreatedAt:        time.Now(),
	},
	CreatedAt: time.Now(),
}

var dPoolData = dmodels.PoolData{
	ID:                uuid.Must(uuid.FromString("285a5c69-74b5-4aec-a0f8-57146dbc3198")),
	PoolID:            uuid.Must(uuid.FromString("6aeb256a-55c9-450c-94b8-e3029eab0ed3")),
	Epoch:             5800,
	ActiveStake:       456215,
	TotalTokensSupply: 10,
	TotalLamports:     20,
	APY:               decimal.Decimal{},
	UnstakeLiquidity:  30,
	DepossitFee:       decimal.Decimal{},
	WithdrawalFee:     decimal.Decimal{},
	RewardsFee:        decimal.Decimal{},
	UpdatedAt:         time.Time{},
	CreatedAt:         time.Time{},
	Pool:              dmodels.Pool{},
}

var poolVD = []*dmodels.PoolValidatorData{
	{
		ID:          uuid.Must(uuid.FromString("9f02b9bd-f065-401f-bb7f-d7048ea7067b")),
		PoolDataID:  uuid.Must(uuid.FromString("6aeb256a-55c9-450c-94b8-e3029eab0ed3")),
		ValidatorID: "id1",
		ActiveStake: 854684,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		Validator:   dmodels.Validator{},
		PoolData:    dmodels.PoolData{},
	},
}

var dValView = dmodels.ValidatorView{
	ID:              "id1",
	Image:           "img1",
	Name:            "val1",
	Delinquent:      true,
	VotePK:          "pk1",
	APY:             decimal.Decimal{},
	StakingAccounts: 500,
	ActiveStake:     100,
	Fee:             decimal.Decimal{},
	Score:           5698,
	SkippedSlots:    decimal.Decimal{},
	DataCenter:      "dc",
	CreatedAt:       time.Time{},
	UpdatedAt:       time.Time{},
}

func TestGetPool(t *testing.T) {
	data := map[string]struct {
		DAO  services.Imp
		Data struct {
			name string
		}
		Result *smodels.PoolDetails
		Err    error
	}{
		"first": {
			Data: struct {
				name string
			}{name: "pool1"},
			Result: &smodels.PoolDetails{
				Pool: smodels.Pool{
					Address:          "addr1",
					Name:             "pool1",
					Image:            "img1",
					Currency:         "coin1",
					ActiveStake:      sol.SOL{decimal.NewFromFloat(0.000456215)},
					TokensSupply:     sol.SOL{decimal.NewFromFloat(0.00000001)},
					TotalLamports:    sol.SOL{decimal.NewFromFloat(0.00000002)},
					APY:              decimal.Decimal{},
					AVGSkippedSlots:  decimal.Decimal{},
					AVGScore:         5698,
					StakingAccounts:  500,
					Delinquent:       1,
					UnstakeLiquidity: sol.SOL{decimal.NewFromFloat(0.00000003)},
					DepossitFee:      decimal.Decimal{},
					WithdrawalFee:    decimal.Decimal{},
					RewardsFee:       decimal.Decimal{},
					ValidatorCount:   1,
					CreatedAt:        time.Time{},
				},
				CreatedAt: time.Time{},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != dPool.ID, poolID is %s", poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
						if validatorID != poolVD[0].ValidatorID {
							return nil, fmt.Errorf("validatorID != %s, validatorID is %s", poolVD[0].ValidatorID, validatorID)
						}
						return &dValView, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != dPool.CoinID {
							return nil, fmt.Errorf("id != %s, id is %s", dPool.CoinID, id)
						}
						return coinArr[0], nil
					},
				},
			},
		},
		"second": {
			Data: struct {
				name string
			}{name: "pool1"},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPool: %s", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"third": {
			Data: struct {
				name string
			}{name: "pool1"},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPool(%s): %w", "pool1", postgres.ErrorRecordNotFounded),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return nil, nil
					},
				},
			},
		},
		"forth": {
			Data: struct {
				name string
			}{name: "pool1"},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPoolData: %s", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"fifth": {
			Data: struct {
				name string
			}{name: "pool1"},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPoolValidatorData: %s", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != dPool.ID, poolID is %s", poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"sixth": {
			Data: struct {
				name string
			}{name: "pool1"},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetValidator(%s): %w", poolVD[0].ValidatorID, fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != dPool.ID, poolID is %s", poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"seventh": {
			Data: struct {
				name string
			}{name: "pool1"},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetCoinByID: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != dPool.ID, poolID is %s", poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
						if validatorID != poolVD[0].ValidatorID {
							return nil, fmt.Errorf("validatorID != %s, validatorID is %s", poolVD[0].ValidatorID, validatorID)
						}
						return &dValView, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
	}

	for s, s2 := range data {
		s2.DAO.Cache = cache.New(time.Minute, time.Minute)
		t.Run(s, func(t *testing.T) {
			pool, err := s2.DAO.GetPool(s2.Data.name)
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}

			t.Run(fmt.Sprintf("pool[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, pool, s2.Result)
			})

		})
	}
}

func TestGetPools(t *testing.T) {
	data := map[string]struct {
		DAO  services.Imp
		Data struct {
			name   string
			sort   string
			desc   bool
			limit  uint64
			offset uint64
		}
		Result []*smodels.PoolDetails
		Err    error
	}{
		"first": {
			Data: struct {
				name   string
				sort   string
				desc   bool
				limit  uint64
				offset uint64
			}{name: "pool1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: []*smodels.PoolDetails{
				{
					Pool: smodels.Pool{
						Address:          "poll_adr1",
						Name:             "Pool1",
						Image:            "img",
						Currency:         "coin1",
						ActiveStake:      sol.SOL{decimal.NewFromFloat(0.000456215)},
						TokensSupply:     sol.SOL{decimal.NewFromFloat(0.00000001)},
						TotalLamports:    sol.SOL{decimal.NewFromFloat(0.00000002)},
						APY:              decimal.Decimal{},
						AVGSkippedSlots:  decimal.Decimal{},
						AVGScore:         5698,
						StakingAccounts:  500,
						Delinquent:       1,
						UnstakeLiquidity: sol.SOL{decimal.NewFromFloat(0.00000003)},
						DepossitFee:      decimal.Decimal{},
						WithdrawalFee:    decimal.Decimal{},
						RewardsFee:       decimal.Decimal{},
						ValidatorCount:   1,
						CreatedAt:        time.Time{},
					},
					CreatedAt: time.Time{},
				},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != %s, condition.Condition.Network = %s", postgres.MainNet, condition.Condition.Network)
						}
						if condition.Condition.Name != "pool1" {
							return nil, fmt.Errorf("condition.Condition.Name != pool1, condition.Condition.Name = %s", condition.Condition.Name)
						}
						if condition.Pagination.Limit != 10 {
							return nil, fmt.Errorf("condition.Pagination.Limit != 10, condition.Pagination.Limit = %d", condition.Pagination.Limit)
						}
						if condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("condition.Pagination.Offset != 0, condition.Pagination.Offset = %d", condition.Pagination.Offset)
						}
						if condition.Sort.PoolSort != 1 {
							return nil, fmt.Errorf("condition.Sort.PoolSort != 1, condition.Sort.PoolSort = %d", condition.Sort.PoolSort)
						}
						if condition.Sort.Desc != true {
							return nil, fmt.Errorf("condition.Sort.Desc != true, condition.Sort.Desc = %v", condition.Sort.Desc)
						}
						return []*dmodels.Pool{poolArr[0]}, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != poolArr[0].ID {
							return nil, fmt.Errorf("poolID != %s, poolID is %s", poolArr[0].ID, poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
						if validatorID != poolVD[0].ValidatorID {
							return nil, fmt.Errorf("validatorID != %s, validatorID is %s", poolVD[0].ValidatorID, validatorID)
						}
						return &dValView, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != dPool.CoinID {
							return nil, fmt.Errorf("id != %s, id is %s", dPool.CoinID, id)
						}
						return coinArr[0], nil
					},
					GetPoolCountFunc: func(condition *postgres.Condition) (int64, error) {
						if condition.Network != postgres.MainNet {
							return 0, fmt.Errorf("condition.Network != mainnet, condition.Network = %s", condition.Network)
						}
						if condition.Name != "pool1" {
							return 0, fmt.Errorf("condition.Name != pool1, condition.Name = %s", condition.Name)
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
			}{name: "", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPool: %w", gorm.ErrRecordNotFound),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != %s, condition.Condition.Network = %s", postgres.MainNet, condition.Condition.Network)
						}
						if condition.Condition.Name == "" {
							return nil, gorm.ErrRecordNotFound
						}
						return nil, nil
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
			}{name: "pool1", sort: "pool stake", desc: true, limit: 0, offset: 0},
			Result: nil,
			Err:    nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != %s, condition.Condition.Network = %s", postgres.MainNet, condition.Condition.Network)
						}
						if condition.Condition.Name != "pool1" {
							return nil, fmt.Errorf("condition.Condition.Name != pool1, condition.Condition.Name = %s", condition.Condition.Name)
						}
						if condition.Pagination.Limit != 0 {
							return nil, fmt.Errorf("condition.Pagination.Limit != 0, condition.Pagination.Limit = %d", condition.Pagination.Limit)
						}
						if condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("condition.Pagination.Offset != 0, condition.Pagination.Offset = %d", condition.Pagination.Offset)
						}
						if condition.Sort.PoolSort != 1 {
							return nil, fmt.Errorf("condition.Sort.PoolSort != 1, condition.Sort.PoolSort = %d", condition.Sort.PoolSort)
						}
						if condition.Sort.Desc != true {
							return nil, fmt.Errorf("condition.Sort.Desc != true, condition.Sort.Desc = %v", condition.Sort.Desc)
						}
						return nil, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != poolArr[0].ID {
							return nil, fmt.Errorf("poolID != %s, poolID is %s", poolArr[0].ID, poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
						if validatorID != poolVD[0].ValidatorID {
							return nil, fmt.Errorf("validatorID != %s, validatorID is %s", poolVD[0].ValidatorID, validatorID)
						}
						return &dValView, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != dPool.CoinID {
							return nil, fmt.Errorf("id != %s, id is %s", dPool.CoinID, id)
						}
						return coinArr[0], nil
					},
					GetPoolCountFunc: func(condition *postgres.Condition) (int64, error) {
						if condition.Network != postgres.MainNet {
							return 0, fmt.Errorf("condition.Network != mainnet, condition.Network = %s", condition.Network)
						}
						if condition.Name != "pool1" {
							return 0, fmt.Errorf("condition.Name != pool1, condition.Name = %s", condition.Name)
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
			}{name: "pool1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetLastPoolData: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != %s, condition.Condition.Network = %s", postgres.MainNet, condition.Condition.Network)
						}
						if condition.Condition.Name != "pool1" {
							return nil, fmt.Errorf("condition.Condition.Name != pool1, condition.Condition.Name = %s", condition.Condition.Name)
						}
						if condition.Pagination.Limit != 10 {
							return nil, fmt.Errorf("condition.Pagination.Limit != 10, condition.Pagination.Limit = %d", condition.Pagination.Limit)
						}
						if condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("condition.Pagination.Offset != 0, condition.Pagination.Offset = %d", condition.Pagination.Offset)
						}
						if condition.Sort.PoolSort != 1 {
							return nil, fmt.Errorf("condition.Sort.PoolSort != 1, condition.Sort.PoolSort = %d", condition.Sort.PoolSort)
						}
						if condition.Sort.Desc != true {
							return nil, fmt.Errorf("condition.Sort.Desc != true, condition.Sort.Desc = %v", condition.Sort.Desc)
						}
						return []*dmodels.Pool{poolArr[0]}, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						return nil, fmt.Errorf("some error")
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
			}{name: "pool1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetValidator: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != %s, condition.Condition.Network = %s", postgres.MainNet, condition.Condition.Network)
						}
						if condition.Condition.Name != "pool1" {
							return nil, fmt.Errorf("condition.Condition.Name != pool1, condition.Condition.Name = %s", condition.Condition.Name)
						}
						if condition.Pagination.Limit != 10 {
							return nil, fmt.Errorf("condition.Pagination.Limit != 10, condition.Pagination.Limit = %d", condition.Pagination.Limit)
						}
						if condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("condition.Pagination.Offset != 0, condition.Pagination.Offset = %d", condition.Pagination.Offset)
						}
						if condition.Sort.PoolSort != 1 {
							return nil, fmt.Errorf("condition.Sort.PoolSort != 1, condition.Sort.PoolSort = %d", condition.Sort.PoolSort)
						}
						if condition.Sort.Desc != true {
							return nil, fmt.Errorf("condition.Sort.Desc != true, condition.Sort.Desc = %v", condition.Sort.Desc)
						}
						return []*dmodels.Pool{poolArr[0]}, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != poolArr[0].ID {
							return nil, fmt.Errorf("poolID != %s, poolID is %s", poolArr[0].ID, poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
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
			}{name: "pool1", sort: "pool stake", desc: true, limit: 10, offset: 0},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPoolCount: %w", gorm.ErrRecordNotFound),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != %s, condition.Condition.Network = %s", postgres.MainNet, condition.Condition.Network)
						}
						if condition.Condition.Name != "pool1" {
							return nil, fmt.Errorf("condition.Condition.Name != pool1, condition.Condition.Name = %s", condition.Condition.Name)
						}
						if condition.Pagination.Limit != 10 {
							return nil, fmt.Errorf("condition.Pagination.Limit != 10, condition.Pagination.Limit = %d", condition.Pagination.Limit)
						}
						if condition.Pagination.Offset != 0 {
							return nil, fmt.Errorf("condition.Pagination.Offset != 0, condition.Pagination.Offset = %d", condition.Pagination.Offset)
						}
						if condition.Sort.PoolSort != 1 {
							return nil, fmt.Errorf("condition.Sort.PoolSort != 1, condition.Sort.PoolSort = %d", condition.Sort.PoolSort)
						}
						if condition.Sort.Desc != true {
							return nil, fmt.Errorf("condition.Sort.Desc != true, condition.Sort.Desc = %v", condition.Sort.Desc)
						}
						return []*dmodels.Pool{poolArr[0]}, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != poolArr[0].ID {
							return nil, fmt.Errorf("poolID != %s, poolID is %s", poolArr[0].ID, poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
						if validatorID != poolVD[0].ValidatorID {
							return nil, fmt.Errorf("validatorID != %s, validatorID is %s", poolVD[0].ValidatorID, validatorID)
						}
						return &dValView, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != dPool.CoinID {
							return nil, fmt.Errorf("id != %s, id is %s", dPool.CoinID, id)
						}
						return coinArr[0], nil
					},
					GetPoolCountFunc: func(condition *postgres.Condition) (int64, error) {
						return 0, gorm.ErrRecordNotFound
					},
				},
			},
		},
	}
	for s, s2 := range data {
		t.Run(s, func(t *testing.T) {
			pools, count, err := s2.DAO.GetPools(s2.Data.name, s2.Data.sort, s2.Data.desc, s2.Data.limit, s2.Data.offset)
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			assert.Equal(t, uint64(len(s2.Result)), count)
			assert.Equal(t, uint64(len(pools)), count)
			for i, pool := range pools {
				t.Run(fmt.Sprintf("governances[%d]", i), func(t *testing.T) {
					assert.DeepEqual(t, pool, s2.Result[i])
				})
			}
		})
	}
}

func TestGetPoolsCurrentStatistic(t *testing.T) {
	data := map[string]struct {
		DAO    services.Imp
		Result *smodels.Statistic
		Err    error
	}{
		"first": {
			Result: &smodels.Statistic{
				Pools:            1,
				ActiveStake:      sol.SOL{decimal.NewFromFloat(0.000456215)},
				TotalSupply:      sol.SOL{decimal.NewFromFloat(0.00000001)},
				AVGSkippedSlots:  decimal.Decimal{},
				MAXPoolsApy:      decimal.Decimal{},
				MAXScore:         5698,
				AVGScore:         5698,
				MINScore:         5698,
				Delinquent:       1,
				UnstakeLiquidity: sol.SOL{decimal.NewFromFloat(0.00000003)},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != mainnet, condition.Condition.Network = %s", condition.Condition.Network)
						}
						return poolArr[:1], nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != poolArr[0].ID {
							return nil, fmt.Errorf("poolID != %s, poolID is %s", poolArr[0].ID, poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
						if validatorID != poolVD[0].ValidatorID {
							return nil, fmt.Errorf("validatorID != %s, validatorID is %s", poolVD[0].ValidatorID, validatorID)
						}
						return &dValView, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != dPool.CoinID {
							return nil, fmt.Errorf("id != %s, id is %s", dPool.CoinID, id)
						}
						return coinArr[0], nil
					},
				},
			},
		},
		"second": {
			Result: nil,
			Err:    nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != mainnet, condition.Condition.Network = %s", condition.Condition.Network)
						}
						return nil, nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != poolArr[0].ID {
							return nil, fmt.Errorf("poolID != %s, poolID is %s", poolArr[0].ID, poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
						if validatorID != poolVD[0].ValidatorID {
							return nil, fmt.Errorf("validatorID != %s, validatorID is %s", poolVD[0].ValidatorID, validatorID)
						}
						return &dValView, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != dPool.CoinID {
							return nil, fmt.Errorf("id != %s, id is %s", dPool.CoinID, id)
						}
						return coinArr[0], nil
					},
				},
			},
		},
		"third": {
			Result: nil,
			Err:    fmt.Errorf("DAO.GetLastPoolData: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != mainnet, condition.Condition.Network = %s", condition.Condition.Network)
						}
						return poolArr[:1], nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"forth": {
			Result: nil,
			Err:    fmt.Errorf("DAO.GetValidator: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != mainnet, condition.Condition.Network = %s", condition.Condition.Network)
						}
						return poolArr[:1], nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != poolArr[0].ID {
							return nil, fmt.Errorf("poolID != %s, poolID is %s", poolArr[0].ID, poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"fifth": {
			Result: nil,
			Err:    fmt.Errorf("DAO.GetCoinByID: %w", gorm.ErrRecordNotFound),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolsFunc: func(condition *postgres.PoolCondition) ([]*dmodels.Pool, error) {
						if condition.Condition.Network != postgres.MainNet {
							return nil, fmt.Errorf("condition.Condition.Network != mainnet, condition.Condition.Network = %s", condition.Condition.Network)
						}
						return poolArr[:1], nil
					},
					GetLastPoolDataFunc: func(poolID uuid.UUID) (*dmodels.PoolData, error) {
						if poolID != poolArr[0].ID {
							return nil, fmt.Errorf("poolID != %s, poolID is %s", poolArr[0].ID, poolID)
						}
						return &dPoolData, nil
					},
					GetPoolValidatorDataFunc: func(condition *postgres.PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return nil, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.PoolID, condition.PoolDataIDs[0])
						}
						return poolVD, nil
					},
					GetValidatorFunc: func(validatorID string) (*dmodels.ValidatorView, error) {
						if validatorID != poolVD[0].ValidatorID {
							return nil, fmt.Errorf("validatorID != %s, validatorID is %s", poolVD[0].ValidatorID, validatorID)
						}
						return &dValView, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						return nil, gorm.ErrRecordNotFound
					},
				},
			},
		},
	}
	for s, s2 := range data {
		s2.DAO.Cache = cache.New(time.Minute, time.Minute)
		t.Run(s, func(t *testing.T) {
			stat, err := s2.DAO.GetPoolsCurrentStatistic()
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			t.Run(fmt.Sprintf("statistics[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, stat, s2.Result)
			})

		})
	}
}

func TestGetPoolStatistic(t *testing.T) {
	data := map[string]struct {
		DAO  services.Imp
		Data struct {
			name      string
			aggregate string
		}
		Result []*smodels.Pool
		Err    error
	}{
		"first": {
			Data: struct {
				name      string
				aggregate string
			}{name: "pool1", aggregate: "quarter"},
			Result: []*smodels.Pool{
				{
					Address:          "addr1",
					Name:             "pool1",
					Image:            "img1",
					Currency:         "coin1",
					ActiveStake:      sol.SOL{decimal.NewFromFloat(0.000456215)},
					TokensSupply:     sol.SOL{decimal.NewFromFloat(0.00000001)},
					TotalLamports:    sol.SOL{decimal.NewFromFloat(0.00000002)},
					APY:              decimal.Decimal{},
					AVGSkippedSlots:  decimal.Decimal{},
					AVGScore:         0,
					StakingAccounts:  0,
					Delinquent:       0,
					UnstakeLiquidity: sol.SOL{decimal.NewFromFloat(0.00000003)},
					DepossitFee:      decimal.Decimal{},
					WithdrawalFee:    decimal.Decimal{},
					RewardsFee:       decimal.Decimal{},
					ValidatorCount:   1,
					CreatedAt:        time.Time{},
				},
			},
			Err: nil,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetPoolStatisticFunc: func(poolID uuid.UUID, aggregate postgres.Aggregate) ([]*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != %s, poolID = %s", dPool.ID, poolID)
						}
						if aggregate != 2 {
							return nil, fmt.Errorf("aggregate != 2, aggregate = %d", aggregate)
						}
						return []*dmodels.PoolData{&dPoolData}, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != dPool.CoinID {
							return nil, fmt.Errorf("id != %s, id is %s", dPool.CoinID, id)
						}
						return coinArr[0], nil
					},
					GetValidatorDataCountFunc: func(condition *postgres.PoolValidatorDataCondition) (int64, error) {
						if condition.PoolDataIDs[0] != dPoolData.ID {
							return 0, fmt.Errorf("condition.PoolDataIDs[0] != %s, condition.PoolDataIDs[0] is %s", dPoolData.ID, condition.PoolDataIDs[0])
						}
						return 1, nil
					},
				},
			},
		},
		"second": {
			Data: struct {
				name      string
				aggregate string
			}{name: "pool1", aggregate: "quarter"},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetPool(%s): %w", "pool1", postgres.ErrorRecordNotFounded),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name == "pool1" {
							return nil, nil
						}
						return &dPool, nil
					},
				},
			},
		},
		"third": {
			Data: struct {
				name      string
				aggregate string
			}{name: "pool1", aggregate: "quarter"},
			Result: nil,
			Err:    fmt.Errorf("some error"),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetPoolStatisticFunc: func(poolID uuid.UUID, aggregate postgres.Aggregate) ([]*dmodels.PoolData, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"forth": {
			Data: struct {
				name      string
				aggregate string
			}{name: "pool1", aggregate: "quarter"},
			Result: nil,
			Err:    fmt.Errorf("DAO.GetCoinByID: %w", fmt.Errorf("some error")),
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetPoolStatisticFunc: func(poolID uuid.UUID, aggregate postgres.Aggregate) ([]*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != %s, poolID = %s", dPool.ID, poolID)
						}
						if aggregate != 2 {
							return nil, fmt.Errorf("aggregate != 2, aggregate = %d", aggregate)
						}
						return []*dmodels.PoolData{&dPoolData}, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						return nil, fmt.Errorf("some error")
					},
				},
			},
		},
		"fifth": {
			Data: struct {
				name      string
				aggregate string
			}{name: "pool1", aggregate: "quarter"},
			Result: nil,
			Err:    gorm.ErrRecordNotFound,
			DAO: services.Imp{
				DAO: &dao.PostgresMock{
					GetPoolFunc: func(name string) (*dmodels.Pool, error) {
						if name != "pool1" {
							return nil, fmt.Errorf("name != pool1, name is %s", name)
						}
						return &dPool, nil
					},
					GetPoolStatisticFunc: func(poolID uuid.UUID, aggregate postgres.Aggregate) ([]*dmodels.PoolData, error) {
						if poolID != dPool.ID {
							return nil, fmt.Errorf("poolID != %s, poolID = %s", dPool.ID, poolID)
						}
						if aggregate != 2 {
							return nil, fmt.Errorf("aggregate != 2, aggregate = %d", aggregate)
						}
						return []*dmodels.PoolData{&dPoolData}, nil
					},
					GetCoinByIDFunc: func(id uuid.UUID) (*dmodels.Coin, error) {
						if id != dPool.CoinID {
							return nil, fmt.Errorf("id != %s, id is %s", dPool.CoinID, id)
						}
						return coinArr[0], nil
					},
					GetValidatorDataCountFunc: func(condition *postgres.PoolValidatorDataCondition) (int64, error) {
						if condition.PoolDataIDs[0] == dPoolData.ID {
							return 0, gorm.ErrRecordNotFound
						}
						return 1, nil
					},
				},
			},
		},
	}
	for s, s2 := range data {
		t.Run(s, func(t *testing.T) {
			stat, err := s2.DAO.GetPoolStatistic(s2.Data.name, s2.Data.aggregate)
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			t.Run(fmt.Sprintf("statistics[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, stat, s2.Result)
			})
		})
	}
}

func TestGetNetworkAPY(t *testing.T) {
	data := map[string]struct {
		DAO    services.Imp
		Result float64
		Err    error
	}{
		"first": {
			Result: 0,
			Err:    fmt.Errorf("%w: %s", cache.KeyWasNotFound, "apy_key"),
		},
	}

	for s, s2 := range data {
		s2.DAO.Cache = cache.New(time.Minute, time.Minute)
		t.Run(s, func(t *testing.T) {
			napy, err := s2.DAO.GetNetworkAPY()
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			t.Run(fmt.Sprintf("statistics[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, napy, s2.Result)
			})

		})
	}
}

func TestGetUSD(t *testing.T) {
	data := map[string]struct {
		DAO    services.Imp
		Result float64
		Err    error
	}{
		"first": {
			Result: 0,
			Err:    fmt.Errorf("%w: %s", cache.KeyWasNotFound, "price_key"),
		},
	}

	for s, s2 := range data {
		s2.DAO.Cache = cache.New(time.Minute, time.Minute)
		t.Run(s, func(t *testing.T) {
			napy, err := s2.DAO.GetUSD()
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			t.Run(fmt.Sprintf("statistics[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, napy, s2.Result)
			})

		})
	}
}
