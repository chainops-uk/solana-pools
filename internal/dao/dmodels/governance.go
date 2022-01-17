package dmodels

import (
	uuid "github.com/satori/go.uuid"
)

type Governance struct {
	ID                 uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();not null;"`
	Name               string    `gorm:"type:varchar(40)"`
	Symbol             string    `gorm:"type:varchar(40)"`
	VoteURL            string    `gorm:"type:text;default:'null';not null;index:idx_governance_vote_url,unique,where:vote_url != 'null';"`
	WebSiteURL         string    `gorm:"type:text;default:'null';not null;index:idx_governance_web_url,unique,where:web_site_url != 'null';"`
	Image              string    `gorm:"type:text;default:'null';not null;"`
	GeckoKey           string    `gorm:"type:varchar(40);not null;"`
	Blockchain         string    `gorm:"type:varchar(40);not null;"`
	ContractAddress    string    `gorm:"type:varchar(120);not null;index:idx_gov_contract_address,unique;"`
	MaximumTokenSupply float64   `gorm:"type:float8;default:0;not null;"`
	CirculatingSupply  float64   `gorm:"type:float8;default:0;not null;"`
	USD                float64   `gorm:"type:float8;default:0;not null;"`
}
