package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (db *DB) GetPool(name string) (pool dmodels.Pool, err error) {
	err = db.First(&pool, "name = ?", name).Error
	return pool, err
}

func (db *DB) GetPools(cond *Condition) ([]dmodels.Pool, error) {
	var pools []dmodels.Pool
	return pools, withCond(db.DB, cond).Find(&pools).Error
}

func (db *DB) GetPoolCount(cond *Condition) (int64, error) {
	var count int64
	if err := withCond(db.DB, cond).Model(&dmodels.Pool{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DB) GetLastPoolData(PoolID uuid.UUID) (*dmodels.PoolData, error) {
	pool := &dmodels.PoolData{}
	if err := db.Where(`pool_id = ?`, PoolID).Order("created_at desc").First(pool).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return pool, nil
}

func (db *DB) UpdatePoolData(pool *dmodels.PoolData) error {
	return db.Save(pool).Error
}

func withCond(db *gorm.DB, cond *Condition) *gorm.DB {
	if cond == nil {
		return db
	}
	if cond.Name != "" {
		db.Where(`name ilike %?%`, cond.Name)
	}
	if cond.Limit != 0 {
		db.Offset(int(cond.Limit))
	}
	if cond.Offset != 0 {
		db.Offset(int(cond.Offset))
	}

	return db
}
