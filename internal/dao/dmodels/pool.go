package dmodels

import uuid "github.com/satori/go.uuid"

type Pool struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	Name     string    `gorm:"type:varchar(100);not null;"`
	Image    string    `gorm:"type:text"`
	Active   bool      `gorm:"not null"`
	Currency string    `gorm:"type:varchar(40);not null;"`
	Address  string    `gorm:"index;not null;"`
	Network  string    `gorm:"type:varchar(50);not null;"`
}
