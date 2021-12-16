package dmodels

import uuid "github.com/satori/go.uuid"

type Coin struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	Name       string    `gorm:"type:varchar(40);not null;index:idx_coin_name,unique;"`
	GeckoKey   string    `gorm:"type:varchar(40);not null;index:idx_coin_gecko_key,unique;"`
	Address    string    `gorm:"type:varchar(120);not null;index:idx_coin_address,unique;"`
	USD        float64   `gorm:"type:float8;not null;default:0;"`
	ThumbImage string    `gorm:"type:varchar(240);not null;default:'NaN';"`
	SmallImage string    `gorm:"type:varchar(240);not null;default:0;default:'NaN';"`
	LargeImage string    `gorm:"type:varchar(240);not null;default:0;default:'NaN';"`
}
