package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (db *DB) GetCoins(cond *CoinCondition) ([]*dmodels.Coin, error) {
	var coins []*dmodels.Coin
	return coins, withCoinCondition(db.DB, cond).Order("name").Find(&coins).Error
}

func (db *DB) GetCoinsCount(cond *CoinCondition) (int64, error) {
	var count int64
	if err := withCoinCondition(db.DB, cond).Model(&dmodels.Coin{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DB) SaveCoin(coins ...*dmodels.Coin) error {
	return db.DB.Save(coins).Error
}

func sortCoin(db *gorm.DB, sort CoinSortType, desc bool) *gorm.DB {
	switch sort {
	case CoinName:
		db = db.Select("coins.*")
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "coins.name",
					},
					Desc: desc,
				},
			},
		})
	case CoinPrice:
		db = db.Select("coins.*")
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "coins.usd",
					},
					Desc: desc,
				},
			},
		})
	}

	return db
}
