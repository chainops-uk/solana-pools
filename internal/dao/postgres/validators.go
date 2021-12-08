package postgres

import (
	"github.com/dfuse-io/solana-go"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"gorm.io/gorm"
)

func (db *DB) GetValidatorByVotePK(key solana.PublicKey) (*dmodels.Validator, error) {
	validator := &dmodels.Validator{}
	err := db.Where("vote_pk = ?", key.String()).First(validator).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return validator, err
}

func (db *DB) GetValidator(validatorID string) (*dmodels.Validator, error) {
	validator := &dmodels.Validator{}
	err := db.Where("id = ?", validatorID).First(validator).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return validator, err
}

func (db *DB) GetValidators(condition *Condition) ([]*dmodels.Validator, error) {
	validators := make([]*dmodels.Validator, 0)
	return validators, withCond(db.DB, condition).Find(&validators).Error
}

func (db *DB) UpdateValidators(validators ...*dmodels.Validator) error {
	return db.Save(&validators).Error
}
