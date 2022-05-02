package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) GetPoolValidatorData(condition *PoolValidatorDataCondition) ([]*dmodels.PoolValidatorData, error) {
	var vd []*dmodels.PoolValidatorData
	return vd, withPoolValidatorDataCondition(db.DB, condition).Find(&vd).Error
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

func (db *DB) GetValidatorDataCount(condition *PoolValidatorDataCondition) (int64, error) {
	i := int64(0)
	return i, withPoolValidatorDataCondition(db.DB.Model(&dmodels.PoolValidatorData{}), condition).Count(&i).Error
}

func withPoolValidatorDataCondition(db *gorm.DB, condition *PoolValidatorDataCondition) *gorm.DB {
	if condition == nil {
		return db
	}

	db = withCond(db, condition.Condition)
	if len(condition.PoolDataIDs) > 0 {
		db = db.Where("pool_data_id in (?)", condition.PoolDataIDs)
	}
	if len(condition.ValidatorIDs) > 0 {
		db = db.Where("validator_id in (?)", condition.ValidatorIDs)
	}

	if condition.Sort != nil {
		db = db.Joins("join material_validator_data_view as validators on validators.id = pool_validator_data.validator_id").
			Select("pool_validator_data.*")
		return sortValidators(db, condition.Sort.ValidatorDataSort, condition.Sort.Desc)
	}

	return db
}

func sortValidators(db *gorm.DB, sort ValidatorDataSortType, desc bool) *gorm.DB {
	switch sort {
	case ValidatorDataPoolStake:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "pool_validator_data.active_stake",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorDataAPY:
		return db.Clauses(clause.OrderBy{

			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.apy",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorDataStake:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.active_stake",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorDataFee:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.fee",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorDataScore:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.score",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorDataSkippedSlot:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.skipped_slots",
					},
					Desc: desc,
				},
			},
		})
	case ValidatorDataDataCenter:
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators.data_center",
					},
					Desc: desc,
				},
			},
		})
	}

	return db
}
