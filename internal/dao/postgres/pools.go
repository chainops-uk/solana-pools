package postgres

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func (db *DB) GetPool(name string) (*dmodels.Pool, error) {
	var pool *dmodels.Pool
	if err := db.First(&pool, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return pool, nil
}

func (db *DB) GetPools(cond *PoolCondition) ([]dmodels.Pool, error) {
	var pools []dmodels.Pool
	return pools, withPoolCondition(db.DB, cond).Find(&pools).Error
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
	if err := db.DB.Where(`pool_id = ?`, PoolID).Order("created_at desc").First(pool).Error; err != nil {
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

func (db *DB) GetPoolStatistic(PoolID uuid.UUID, aggregate Aggregate) ([]*dmodels.PoolData, error) {
	var data []*dmodels.PoolData
	w, err := aggregateByDate(aggregate, db.DB)
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

func withPoolCondition(db *gorm.DB, condition *PoolCondition) *gorm.DB {
	if condition == nil {
		return db
	}

	if condition.Condition != nil {
		if condition.Condition.Name != "" {
			db = db.Where(`pools.name ilike ?`, "%"+condition.Condition.Name+"%")
		}
		condition.Condition.Name = ""
	}

	db = withCond(db, condition.Condition)

	if condition.Sort != nil {
		return sortPoolData(db, condition.Sort.PoolSort, condition.Sort.Desc)
	}

	return db
}

func sortPoolData(db *gorm.DB, sort PoolDataSortType, desc bool) *gorm.DB {
	db = db.Joins("left join pool_data on pools.id = pool_data.pool_id").
		Where(`pool_data.created_at = (SELECT max(t1.created_at) FROM pool_data t1 WHERE  t1.pool_id = pools.id)`).
		Group("pools.id")
	switch sort {
	case PoolAPY:
		db = db.Group(`"pool_data"."id"`).Select("pools.*")
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "pool_data.apy",
					},
					Desc: desc,
				},
			},
		})
	case PoolStake:
		db = db.Group(`"pool_data"."id"`).Select("pools.*")
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "pool_data.active_stake",
					},
					Desc: desc,
				},
			},
		})
	case PoolValidators:
		db = db.Joins("left join pool_validator_data on pool_data.id = pool_validator_data.pool_data_id").
			Select("pools.*, count(pool_validator_data.*) as validators")
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "validators",
					},
					Desc: desc,
				},
			},
		})
	case PoolScore:
		db = db.Joins("left join pool_validator_data on pool_data.id = pool_validator_data.pool_data_id").
			Joins("join validators on pool_validator_data.validator_id = validators.id").
			Select("pools.*, avg(validators.score) as avg_score")
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "avg_score",
					},
					Desc: desc,
				},
			},
		})
	case PoolSkippedSlot:
		db = db.Joins("left join pool_validator_data on pool_data.id = pool_validator_data.pool_data_id").
			Joins("join validators on pool_validator_data.validator_id = validators.id").
			Select("pools.*, avg(validators.skipped_slots) as avg_skipped_slots")
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "avg_skipped_slots",
					},
					Desc: desc,
				},
			},
		})
	case PoolTokenPrice:
		db = db.Select("pools.*, (CASE WHEN pool_data.total_tokens_supply IS NULL THEN 0 WHEN pool_data.total_tokens_supply = 0 THEN 0 ELSE pool_data.total_lamports / pool_data.total_tokens_supply END) as price").
			Group("pool_data.id")
		return db.Clauses(clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column: clause.Column{
						Name: "price",
					},
					Desc: desc,
				},
			},
		})
	}

	return db
}

func aggregateByDate(aggregate Aggregate, db *gorm.DB) (*gorm.DB, error) {
	switch aggregate {
	case Week:
		return db.Where(`"created_at"::date between ? AND ?`, time.Now().AddDate(0, 0, -7), time.Now()).
			Where(`created_at = (SELECT max(t1.created_at) FROM pool_data t1 WHERE  t1.pool_id = "pool_data".pool_id and t1.created_at::date = pool_data.created_at::date)`), nil
	case Month:
		return db.Where(`"created_at"::date between ? AND ?`, time.Now().AddDate(0, -1, 0), time.Now()).
			Where(`created_at = (SELECT max(t1.created_at) FROM pool_data t1 WHERE  t1.pool_id = "pool_data".pool_id and t1.created_at::date = pool_data.created_at::date)`), nil
	case Quarter:
		return db.Where(`"created_at"::date between ? AND ?`, time.Now().AddDate(0, -3, 0), time.Now()).
			Where(`"created_at" = (SELECT max(t1.created_at) FROM pool_data t1 WHERE  t1.pool_id = "pool_data".pool_id and date_part('year', "pool_data"."created_at") = date_part('year', t1."created_at") and date_part('week', "pool_data"."created_at") = date_part('week', t1."created_at"))`), nil
	case HalfYear:
		return db.Where(`"created_at"::date between ? AND ?`, time.Now().AddDate(0, -3, 0), time.Now()).
			Where(`"created_at" = (SELECT max(t1.created_at) FROM pool_data t1 WHERE  t1.pool_id = "pool_data".pool_id and date_part('year', "pool_data"."created_at") = date_part('year', t1."created_at") and date_part('week', "pool_data"."created_at") = date_part('week', t1."created_at"))`), nil
	case Year:
		return db.Where(`"created_at"::date between ? AND ?`, time.Now().AddDate(-1, 0, 0), time.Now()).
			Where(`"created_at" = (SELECT max(t1.created_at) FROM pool_data t1 WHERE  t1.pool_id = "pool_data".pool_id and date_part('year', "pool_data"."created_at") = date_part('year', t1."created_at") and date_part('month', "pool_data"."created_at") = date_part('month', t1."created_at"))`), nil
	default:
		return nil, nil
	}
}

func (db *DB) UpdatePoolData(pool *dmodels.PoolData) error {
	return db.Save(pool).Error
}

/*
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
*/
