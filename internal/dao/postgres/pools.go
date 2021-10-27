package postgres

import "github.com/everstake/solana-pools/internal/dao/dmodels"

func (db *DB) GetPool(name string) (pool dmodels.Pool, err error) {
	err = db.First(&pool, "name = ?", name).Error
	return pool, err
}

func (db *DB) GetPools() (pools []dmodels.Pool, err error) {
	err = db.Find(&pools).Error
	return pools, err
}

func (db *DB) UpdatePool(pool dmodels.Pool) error {
	return db.Save(pool).Error
}
