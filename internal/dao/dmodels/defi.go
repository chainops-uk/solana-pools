package dmodels

import (
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type DEFI struct {
	ID              uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	LiquidityPoolID uuid.UUID       `gorm:"type:uuid;not null;"`
	SaleCoinID      uuid.UUID       `gorm:"type:uuid;not null;"`
	BuyCoinID       uuid.UUID       `gorm:"type:uuid;not null;"`
	Liquidity       float64         `gorm:"type:float8;not null;"`
	APY             decimal.Decimal `gorm:"type:decimal(24,9);not null;"`
	SaleCoin        Coin            `gorm:"foreignKey:SaleCoinID;constraint:OnUpdate:CASCADE,OnDelete:Restrict;"`
	BuyCoin         Coin            `gorm:"foreignKey:BuyCoinID;constraint:OnUpdate:CASCADE,OnDelete:Restrict;"`
	LiquidityPool   LiquidityPool   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:Restrict;"`
}
