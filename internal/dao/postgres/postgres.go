package postgres

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/pkg/logger/zapgorm"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

var autoMigrateModels = []interface{}{
	&dmodels.Pool{},
	&dmodels.PoolData{},
	&dmodels.Validator{},
	&dmodels.PoolValidatorData{},
}

func NewDB(dsn string) (db *DB, err error) {
	z, _ := zap.NewProduction()
	logger := zapgorm.New(z)
	logger.SetAsDefault()
	logger.LogMode(gormlogger.Error)
	logger.IgnoreRecordNotFoundError = true
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return db, fmt.Errorf("gorm.Open: %s", err.Error())
	}
	err = d.AutoMigrate(autoMigrateModels...)
	if err != nil {
		return db, fmt.Errorf("gorm.AutoMigrate: %s", err.Error())
	}
	return &DB{d}, nil
}
