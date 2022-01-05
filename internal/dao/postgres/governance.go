package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db *DB) GetGovernance(cond *GovernanceCondition) ([]*dmodels.Governance, error) {
	var gov []*dmodels.Governance
	return gov, withGovernanceCondition(db.DB, cond).Order("name").Find(&gov).Error
}

func (db *DB) GetGovernanceCount(cond *GovernanceCondition) (int64, error) {
	var count int64
	if err := withGovernanceCondition(db.DB, cond).Model(&dmodels.Governance{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DB) SaveGovernance(gov ...*dmodels.Governance) error {
	return db.Save(gov).Error
}

func withGovernanceCondition(db *gorm.DB, cond *GovernanceCondition) *gorm.DB {
	if cond == nil {
		return db
	}
	db = withCond(db, cond.Condition)

	if cond.Sort != nil {
		db = sortGovernance(db, cond.Sort.Sort, cond.Sort.Desc)
	}

	return db
}

func sortGovernance(db *gorm.DB, sort GovernanceSortType, desc bool) *gorm.DB {

	switch sort {
	case GovernanceName:
		db = db.Select("governances.*")
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "governances.name",
					},
					Desc: desc,
				},
			},
		})
	case GovernancePrice:
		db = db.Select("governances.*")
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "governances.usd",
					},
					Desc: desc,
				},
			},
		})
	}

	return db
}
