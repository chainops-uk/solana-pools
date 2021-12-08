package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (db *DB) GetPoolValidatorData(poolID uuid.UUID) (vd []*dmodels.PoolValidatorData, err error) {
	err = db.Where("pool_data_id = ?", poolID).Find(&vd).Error
	return vd, err
}

func (db *DB) CreatePoolValidatorData(validatorsPoolData ...*dmodels.PoolValidatorData) error {
	if len(validatorsPoolData) == 0 {
		return nil
	}
	return db.Create(&validatorsPoolData).Error
}

func (db *DB) DeleteValidators(poolID uuid.UUID) error {
	return db.Where("pool_data_id = ?", poolID).Delete(&dmodels.PoolValidatorData{}).Error
}

func (db *DB) GetValidatorCount(condition *PoolValidatorDataCondition) (int64, error) {
	i := int64(0)
	return i, withConditionPoolValidatorData(db.DB.Model(&dmodels.PoolValidatorData{}), condition).Count(&i).Error
}

func withConditionPoolValidatorData(db *gorm.DB, condition *PoolValidatorDataCondition) *gorm.DB {
	if condition == nil {
		return db
	}
	db = withCond(db, &condition.Condition)
	if len(condition.PoolDataIDs) > 0 {
		db = db.Where("pool_data_id in ?", condition.PoolDataIDs)
	}
	if len(condition.ValidatorIDs) > 0 {
		db = db.Where("validator_id in ?", condition.ValidatorIDs)
	}

	return db
}
