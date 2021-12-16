package dmodels

import (
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type Governance struct {
	ID                    uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	Name                  string          `gorm:"type:varchar(40)"`
	Image                 string          `gorm:"type:text;default:'null';not null;"`
	CoinName              string          `gorm:"type:varchar(40);not null;"`
	Blockchain            string          `gorm:"type:varchar(40);not null;"`
	ContractAddress       string          `gorm:"type:varchar(120);not null;index:idx_gov_contract_address,unique;"`
	VoteURL               string          `gorm:"type:text;default:'null';not null;"`
	About                 string          `gorm:"type:text;default:'null';not null;"`
	Vote                  string          `gorm:"type:text;default:'null';not null;"`
	Trade                 string          `gorm:"type:text;default:'null';not null;"`
	Exchange              string          `gorm:"type:text;default:'null';not null;"`
	MaximumTokenSupply    float64         `gorm:"type:float8;not null;"`
	CirculatingSupply     float64         `gorm:"type:float8;not null;"`
	USD                   float64         `gorm:"type:float8;not null;"`
	DAOTreasury           decimal.Decimal `gorm:"type:numeric(5,2);not null;"`
	Investors             decimal.Decimal `gorm:"type:numeric(5,2);not null;"`
	InitialLidoDevelopers decimal.Decimal `gorm:"type:numeric(5,2);not null;"`
	Foundation            decimal.Decimal `gorm:"type:numeric(5,2);not null;"`
	Validators            decimal.Decimal `gorm:"type:numeric(5,2);not null;"`
}
