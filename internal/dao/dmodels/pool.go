package dmodels

import uuid "github.com/satori/go.uuid"

type Pool struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	Name    string    `gorm:"type:varchar(100);not null;"`
	Active  bool      `gorm:"not null"`
	CoinID  uuid.UUID `gorm:"type:uuid;not null;"`
	Address string    `gorm:"index;not null;"`
	Network string    `gorm:"type:varchar(50);not null;"`
	Coin    Coin      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:Restrict;"`
}
