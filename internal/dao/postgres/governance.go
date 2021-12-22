package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
)

func (db *DB) GetGovernance(cond *CoinCondition) ([]*dmodels.Governance, error) {
	var gov []*dmodels.Governance
	return gov, withCoinCondition(db.DB, cond).Order("name").Find(&gov).Error
}

func (db *DB) GetGovernanceCount(cond *CoinCondition) (int64, error) {
	var count int64
	if err := withCoinCondition(db.DB, cond).Model(&dmodels.Governance{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DB) SaveGovernance(gov ...*dmodels.Governance) error {
	return db.Save(gov).Error
}

func (db *DB) withGovernanceCondition(cond *CoinCondition) (int64, error) {
	var count int64
	if err := withCoinCondition(db.DB, cond).Model(&dmodels.Governance{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
