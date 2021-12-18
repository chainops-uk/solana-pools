package dmodels

import uuid "github.com/satori/go.uuid"

type LiquidityPool struct {
	ID    uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	Name  string    `gorm:"type:varchar(40)"`
	Image string    `gorm:"type:text;default:'null';not null;"`
	URL   string    `gorm:"type:text;default:'null';not null;"`
}
