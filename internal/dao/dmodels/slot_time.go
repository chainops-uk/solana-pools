package dmodels

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type SlotTime struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	SlotTime  float64   `gorm:"type:float8;default:0;not null;"`
	Epoch     uint64    `gorm:"type:int8;default:0;not null;"`
	CreatedAt time.Time `gorm:"index;not null"`
}
