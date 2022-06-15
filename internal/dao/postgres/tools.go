package postgres

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Network string

var ErrorRecordNotFounded = errors.New("record not founded")

const (
	TestNet = "testnet"
	MainNet = "mainnet"
)

type ValidatorCondition struct {
	*Condition
	Epochs       []uint64
	PoolDataIDs  []uuid.UUID
	ValidatorIDs []string
	Sort         *ValidatorSort
}

type PoolValidatorDataCondition struct {
	*Condition
	PoolDataIDs  []uuid.UUID
	ValidatorIDs []string
	Sort         *ValidatorDataSort
}

type CoinCondition struct {
	*Condition
	GeckoIDs []string
	Names    []string
	Name     string
	CoinSort *CoinSort
}

type PoolCondition struct {
	*Condition
	Sort *PoolDataSort
}

type Condition struct {
	IDs     []uuid.UUID
	Names   []string
	Name    string
	Epochs  []uint64
	Network Network
	Pagination
}

type PoolDataSort struct {
	Epoch    uint64
	PoolSort PoolDataSortType
	Desc     bool
}

type PoolDataSortType int

const (
	PoolAPY = PoolDataSortType(iota)
	PoolStake
	PoolValidators
	PoolScore
	PoolSkippedSlot
	PoolTokenPrice
)

func SearchPoolSort(sort string) PoolDataSortType {
	switch sort {
	case "pool stake":
		return PoolStake
	case "validators":
		return PoolValidators
	case "score":
		return PoolScore
	case "skipped slot":
		return PoolSkippedSlot
	case "token price":
		return PoolTokenPrice
	default:
		return PoolAPY
	}
}

type CoinSort struct {
	Sort CoinSortType
	Desc bool
}

type CoinSortType int

const (
	CoinName = CoinSortType(iota)
	CoinPrice
)

func SearchCoinSort(sort string) CoinSortType {
	switch sort {
	case "price":
		return CoinPrice
	default:
		return CoinName
	}
}

type ValidatorDataSortType int

type ValidatorDataSort struct {
	ValidatorDataSort ValidatorDataSortType
	Desc              bool
}

const (
	ValidatorDataAPY = ValidatorDataSortType(iota)
	ValidatorDataPoolStake
	ValidatorDataStake
	ValidatorDataFee
	ValidatorDataScore
	ValidatorDataSkippedSlot
	ValidatorDataDataCenter
)

func SearchValidatorDataSort(sort string) ValidatorDataSortType {
	switch sort {
	case "pool stake":
		return ValidatorDataPoolStake
	case "stake":
		return ValidatorDataStake
	case "fee":
		return ValidatorDataFee
	case "score":
		return ValidatorDataScore
	case "skipped slot":
		return ValidatorDataSkippedSlot
	case "data center":
		return ValidatorDataDataCenter
	default:
		return ValidatorDataAPY
	}
}

type ValidatorSortType int

type ValidatorSort struct {
	ValidatorSort ValidatorSortType
	Desc          bool
}

const (
	ValidatorAPY = ValidatorSortType(iota)
	ValidatorPoolStake
	ValidatorStake
	ValidatorFee
	ValidatorScore
	ValidatorSkippedSlot
	ValidatorDataCenter
	StakingAccounts
)

func SearchValidatorSort(sort string) ValidatorSortType {
	switch sort {
	case "pool stake":
		return ValidatorPoolStake
	case "stake":
		return ValidatorStake
	case "fee":
		return ValidatorFee
	case "score":
		return ValidatorScore
	case "skipped slot":
		return ValidatorSkippedSlot
	case "data center":
		return ValidatorDataCenter
	case "staking accounts":
		return StakingAccounts
	default:
		return ValidatorAPY
	}
}

type Pagination struct {
	Limit  uint64
	Offset uint64
}
type Aggregate int8

const (
	Month = Aggregate(iota)
	Week
	Quarter
	HalfYear
	Year
)

func SearchAggregate(name string) Aggregate {
	switch name {
	case "month":
		return Month
	case "quarter":
		return Quarter
	case "half-year":
		return HalfYear
	case "year":
		return Year
	default:
		return Week
	}
}

type GovernanceCondition struct {
	*Condition
	Sort *GovernanceSort
}

type GovernanceSort struct {
	Sort GovernanceSortType
	Desc bool
}

type GovernanceSortType int

const (
	GovernanceName = GovernanceSortType(iota)
	GovernancePrice
)

func SearchGovernanceSort(sort string) GovernanceSortType {
	switch sort {
	case "price":
		return GovernancePrice
	default:
		return GovernanceName
	}
}

func withCond(db *gorm.DB, cond *Condition) *gorm.DB {
	if cond == nil {
		return db
	}
	if len(cond.IDs) > 0 {
		db = db.Where(`id IN (?)`, cond.IDs)
	}
	if len(cond.Names) > 0 {
		db = db.Where(`name IN (?)`, cond.Names)
	}
	if cond.Name != "" {
		db = db.Where(`name ilike ?`, "%"+cond.Name+"%")
	}
	switch cond.Network {
	case MainNet, TestNet:
		db = db.Where(`network = ?`, string(cond.Network))
	}
	if cond.Limit > 0 {
		db = db.Limit(int(cond.Limit))
	}
	if cond.Offset > 0 {
		db = db.Offset(int(cond.Offset))
	}

	return db
}

func withCoinCondition(db *gorm.DB, cond *CoinCondition) *gorm.DB {
	if cond == nil {
		return db
	}
	db = withCond(db, cond.Condition)
	if len(cond.GeckoIDs) > 0 {
		db = db.Where(`gecko_key IN (?)`, cond.GeckoIDs)
	}
	if len(cond.Names) > 0 {
		db = db.Where(`name IN (?)`, cond.Names)
	}
	if cond.Name != "" {
		db = db.Where(`name ilike ?`, "%"+cond.Name+"%")
	}

	if cond.CoinSort != nil {
		db = sortCoin(db, cond.CoinSort.Sort, cond.CoinSort.Desc)
	}

	return db
}
