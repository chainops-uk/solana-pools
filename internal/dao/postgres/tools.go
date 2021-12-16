package postgres

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Network string

const (
	TestNet = "testnet"
	MainNet = "mainnet"
)

type PoolValidatorDataCondition struct {
	Condition
	PoolDataIDs  []uuid.UUID
	ValidatorIDs []string
}

type Condition struct {
	IDs     []uuid.UUID
	Name    string
	Epoch   []uint64
	Network Network
	Pagination
}
type Pagination struct {
	Limit  uint64
	Offset uint64
}
type Aggregate int8

const (
	Month = Aggregate(iota)
	Week
	Year
)

func SearchAggregate(name string) Aggregate {
	switch name {
	case "month":
		return Month
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
