package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"gorm.io/gorm"
)

func (db *DB) GetLiquidityPools(cond *Condition) ([]*dmodels.LiquidityPool, error) {
	var lp []*dmodels.LiquidityPool
	return lp, withCond(db.DB, cond).Order("name").Find(&lp).Error
}

func (db *DB) GetLiquidityPool(cond *Condition) (*dmodels.LiquidityPool, error) {
	var pool *dmodels.LiquidityPool
	if err := withCond(db.DB, cond).First(&pool).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return pool, nil
}
