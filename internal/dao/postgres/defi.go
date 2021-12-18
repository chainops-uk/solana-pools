package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DeFiCondition struct {
	Condition
	LiquidityPoolIDs []uuid.UUID
	SaleCoinID       []uuid.UUID
	BuyCoinID        []uuid.UUID
}

func (db *DB) GetDEFIs(cond *DeFiCondition) ([]*dmodels.DEFI, error) {
	var lp []*dmodels.DEFI
	return lp, withDeFiCondition(db.DB, cond).Find(&lp).Error
}

func (db *DB) DeleteDeFis(cond *DeFiCondition) error {
	return withDeFiCondition(db.DB, cond).Delete(&dmodels.DEFI{}).Error
}

func (db *DB) SaveDEFIs(defiData ...*dmodels.DEFI) error {
	if len(defiData) == 0 {
		return nil
	}
	return db.Save(&defiData).Error
}

func withDeFiCondition(db *gorm.DB, cond *DeFiCondition) *gorm.DB {
	if cond == nil {
		return db
	}
	db = withCond(db, &cond.Condition)
	if len(cond.LiquidityPoolIDs) > 0 {
		db = db.Where("liquidity_pool_id in (?)", cond.LiquidityPoolIDs)
	}
	if len(cond.SaleCoinID) > 0 {
		db = db.Where("sale_coin_id in (?)", cond.SaleCoinID)
	}
	if len(cond.BuyCoinID) > 0 {
		db = db.Where("buy_coin_id in (?)", cond.BuyCoinID)
	}

	return db
}
