package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (db *DB) GetCoinByID(id uuid.UUID) (*dmodels.Coin, error) {
	coin := &dmodels.Coin{}
	if err := db.First(&coin, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
	}
	return coin, nil
}

func (db *DB) GetCoins(cond *Condition) ([]*dmodels.Coin, error) {
	var coins []*dmodels.Coin
	return coins, withCond(db.DB, cond).Order("name").Find(&coins).Error
}

func (db *DB) GetCoinsCount(cond *Condition) (int64, error) {
	var count int64
	if err := withCond(db.DB, cond).Model(&dmodels.Coin{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DB) SaveCoin(coins ...*dmodels.Coin) error {
	return db.DB.Save(coins).Error
}
