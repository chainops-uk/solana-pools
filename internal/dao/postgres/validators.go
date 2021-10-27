package postgres

import "github.com/everstake/solana-pools/internal/dao/dmodels"

func (db *DB) GetValidators(poolID uint64) (pools []dmodels.Validator, err error) {
	err = db.Where("pool_id = ?", poolID).Find(&pools).Error
	return pools, err
}

func (db *DB) CreateValidators(validators []dmodels.Validator) error {
	return db.Create(&validators).Error
}

func (db *DB) DeleteValidators(poolID uint64) error {
	return db.Where("pool_id = ?", poolID).Delete(&dmodels.Validator{}).Error
}
