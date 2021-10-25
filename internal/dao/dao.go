package dao

import (
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
)

type (
	DAO interface {
		Postgres
	}
	Postgres interface {
		// pools
		GetPool(name string) (pool dmodels.Pool, err error)
		GetPools() (pools []dmodels.Pool, err error)
		UpdatePool(pool dmodels.Pool) error

		// validators
		GetValidator(name string) (v dmodels.Validator, err error)
		GetValidators() (pools []dmodels.Validator, err error)
		CreateValidators(pools []dmodels.Validator) error
		DeleteValidators(poolID uint64) error
	}
	Imp struct {
	}
)

func NewDAO(cfg config.Env) (d DAO, err error) {
	return Imp{}, nil
}

// todo delete & implement
func (i Imp) GetPool(name string) (pool dmodels.Pool, err error) {
	panic("implement me")
}

func (i Imp) GetPools() (pools []dmodels.Pool, err error) {
	panic("implement me")
}

func (i Imp) UpdatePool(pool dmodels.Pool) error {
	panic("implement me")
}

func (i Imp) GetValidator(name string) (v dmodels.Validator, err error) {
	panic("implement me")
}

func (i Imp) GetValidators() (pools []dmodels.Validator, err error) {
	panic("implement me")
}

func (i Imp) UpdateValidator(pool dmodels.Validator) error {
	panic("implement me")
}

func (i Imp) CreateValidators(pools []dmodels.Validator) error {
	panic("implement me")
}

func (i Imp) DeleteValidators(poolID uint64) error {
	panic("implement me")
}
