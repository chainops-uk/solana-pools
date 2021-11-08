package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	uuid "github.com/satori/go.uuid"
)

func (db *DB) GetPool(name string) (pool dmodels.Pool, err error) {
	err = db.First(&pool, "name = ?", name).Error
	return pool, err
}

func (db *DB) GetPools() (pools []dmodels.Pool, err error) {
	err = db.Find(&pools).Error
	return pools, err
}

func (db *DB) GetLastPoolData(PoolID uuid.UUID) (pool *dmodels.PoolData, err error) {
	err = db.Where(`pool_id = ?`, PoolID).Order("create_at desc").Limit(1).Find(pool).Error
	return pool, err
}

func (db *DB) UpdatePoolData(pool *dmodels.PoolData) error {
	return db.Save(pool).Error
}
