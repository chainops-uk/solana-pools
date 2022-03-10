package postgres

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/pkg/logger/zapgorm"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	gormlogger_gorm "gorm.io/driver/postgres"
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
	&dmodels.ValidatorData{},
	&dmodels.PoolValidatorData{},
	&dmodels.Coin{},
	&dmodels.Governance{},
	&dmodels.LiquidityPool{},
	&dmodels.DEFI{},
}

func NewDB(dsn string) (db *DB, err error) {
	z, _ := zap.NewProduction()
	logger := zapgorm.New(z)
	logger.SetAsDefault()
	logger.IgnoreRecordNotFoundError = true
	d, err := gorm.Open(gormlogger_gorm.Open(dsn), &gorm.Config{
		Logger: logger.LogMode(gormlogger.Error),
	})
	if err != nil {
		return db, fmt.Errorf("gorm.Open: %s", err.Error())
	}
	err = d.AutoMigrate(autoMigrateModels...)
	if err != nil {
		return db, fmt.Errorf("gorm.AutoMigrate: %s", err.Error())
	}

	dbm, err := d.DB()
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(dbm, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	fsrc, err := (&file.File{}).Open("file://migrations")
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance(
		"file",
		fsrc,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return &DB{d}, nil
}

func migrationUP() {

}
