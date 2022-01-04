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

type PoolValidatorDataCondition struct {
	*Condition
	PoolDataIDs  []uuid.UUID
	ValidatorIDs []string
	Sort         *ValidatorSort
}

type CoinCondition struct {
	*Condition
	GeckoIDs []string
}

type PoolCondition struct {
	*Condition
	Sort *PoolSort
}

type Condition struct {
	IDs     []uuid.UUID
	Names   []string
	Name    string
	Epoch   []uint64
	Network Network
	Pagination
}

type PoolSort struct {
	PoolSort PoolSortType
	Desc     bool
}

type PoolSortType int

const (
	PoolAPY = PoolSortType(iota)
	PoolStake
	PoolValidators
	PoolScore
	PoolSkippedSlot
	PoolTokenPrice
)

func SearchPoolSort(sort string) PoolSortType {
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
	return db
}
