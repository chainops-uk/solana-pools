package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
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

func (db *DB) GetLastPoolData(PoolID uuid.UUID, condition *Condition) (*dmodels.PoolData, error) {
	pool := &dmodels.PoolData{}
	if err := db.Where(`pool_id = ?`, PoolID).Order("created_at desc").First(pool).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return pool, nil
}

func (db *DB) GetLastEpochPoolData(PoolID uuid.UUID, currentEpoch uint64) (*dmodels.PoolData, error) {
	pool := &dmodels.PoolData{}
	if err := db.Where(`pool_id = ?`, PoolID).
		Where(`epoch < ?`, currentEpoch).
		Order("created_at desc").Limit(1).First(pool).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return pool, nil
}

func (db *DB) GetPoolStatistic(PoolID uuid.UUID, aggregate Aggregate, from, to time.Time) ([]*dmodels.PoolData, error) {
	var data []*dmodels.PoolData
	w, err := aggregateByDate(aggregate, from, to, db.DB)
	if err != nil {
		return nil, err
	}
	if err := w.Where(`pool_id = ?`, PoolID).
		Order("created_at").Find(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

func aggregateByDate(aggregate Aggregate, from, to time.Time, db *gorm.DB) (*gorm.DB, error) {
	switch aggregate {
	case Day:
		return db.Where(`"created_at"::date between ? AND ?`, from, to).
			Where(`created_at = (SELECT max(t1.created_at) FROM pool_data t1 WHERE  t1.pool_id = "pool_data".pool_id and t1.created_at::date = pool_data.created_at::date)`).Distinct(), nil
	case Month:
		return db.Where(`"created_at"::date between ? AND ?`, from, to).
			Where(`"created_at" = (SELECT max(t1.created_at) FROM pool_data t1 WHERE  t1.pool_id = "pool_data".pool_id and date_part('year', "pool_data"."created_at") = date_part('year', t1."created_at") and date_part('week', "pool_data"."created_at") = date_part('week', t1."created_at"))`), nil
	case Week:
		return db.Where(`"created_at"::date between ? AND ?`, from, to).
			Where(`"created_at" = (SELECT max(t1.created_at) FROM pool_data t1 WHERE  t1.pool_id = "pool_data".pool_id and date_part('year', "pool_data"."created_at") = date_part('year', t1."created_at") and date_part('month', "pool_data"."created_at") = date_part('month', t1."created_at"))`), nil
	case Year:
		return db.Where(`"created_at"::date between ? AND ?`, from, to).
			Where(`"created_at" = (SELECT max(t1.created_at) FROM pool_data t1 WHERE  t1.pool_id = "pool_data".pool_id and date_part('year', "pool_data"."created_at") = date_part('year', t1."created_at"))`), nil
	default:
		return nil, nil
	}
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
