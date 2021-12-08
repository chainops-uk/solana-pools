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
	Day = Aggregate(iota)
	Month
	Week
	Year
)

func SearchAggregate(name string) Aggregate {
	switch name {
	case "month":
		return Month
	case "week":
		return Week
	case "year":
		return Year
	default:
		return Day
	}
}

func withCond(db *gorm.DB, cond *Condition) *gorm.DB {
	if cond == nil {
		return db
	}
	if cond.Name != "" {
		db.Where(`name ilike %?%`, cond.Name)
	}
	switch cond.Network {
	case MainNet, TestNet:
		db.Where(`network = ?`, cond.Network)
	}
	if cond.Limit != 0 {
		db.Offset(int(cond.Limit))
	}
	if cond.Offset != 0 {
		db.Offset(int(cond.Offset))
	}

	return db
}
