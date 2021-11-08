package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	uuid "github.com/satori/go.uuid"
)

func (db *DB) GetValidators(poolID uuid.UUID) (pools []*dmodels.Validator, err error) {
	err = db.Where("pool_id = ?", poolID).Find(&pools).Error
	return pools, err
}

func (db *DB) CreateValidator(validators ...*dmodels.Validator) error {
	return db.Create(&validators).Error
}

func (db *DB) DeleteValidators(poolID uuid.UUID) error {
	return db.Where("pool_id = ?", poolID).Delete(&dmodels.Validator{}).Error
}
