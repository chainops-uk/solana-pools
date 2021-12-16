package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
)

func (db *DB) GetGovernance(cond *Condition) ([]*dmodels.Governance, error) {
	var gov []*dmodels.Governance
	return gov, withCond(db.DB, cond).Order("name").Find(&gov).Error
}

func (db *DB) GetGovernanceCount(cond *Condition) (int64, error) {
	var count int64
	if err := withCond(db.DB, cond).Model(&dmodels.Governance{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DB) SaveGovernance(gov ...*dmodels.Governance) error {
	return db.Save(gov).Error
}
