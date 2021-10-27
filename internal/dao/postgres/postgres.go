package postgres

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

var autoMigrateModels = []interface{}{
	&dmodels.Pool{},
	&dmodels.Validator{},
}

func NewDB(dsn string) (db *DB, err error) {
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("gorm.Open: %s", err.Error())
	}
	err = d.AutoMigrate(autoMigrateModels...)
	if err != nil {
		return db, fmt.Errorf("gorm.AutoMigrate: %s", err.Error())
	}
	return &DB{d}, nil
}
