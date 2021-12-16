package smodels

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/shopspring/decimal"
)

type Governance struct {
	Name                  string
	Image                 string
	CoinName              string
	Blockchain            string
	ContractAddress       string
	VoteURL               string
	About                 string
	Vote                  string
	Trade                 string
	Exchange              string
	MaximumTokenSupply    float64
	CirculatingSupply     float64
	USD                   float64
	DAOTreasury           decimal.Decimal
	Investors             decimal.Decimal
	InitialLidoDevelopers decimal.Decimal
	Foundation            decimal.Decimal
	Validators            decimal.Decimal
}

func (g *Governance) Set(governance *dmodels.Governance) *Governance {
	g.Name = governance.Name
	g.Image = governance.Image
	g.CoinName = governance.CoinName
	g.Blockchain = governance.Blockchain
	g.ContractAddress = governance.ContractAddress
	g.VoteURL = governance.VoteURL
	g.About = governance.About
	g.Vote = governance.Vote
	g.Trade = governance.Trade
	g.Exchange = governance.Exchange
	g.MaximumTokenSupply = governance.MaximumTokenSupply
	g.CirculatingSupply = governance.CirculatingSupply
	g.USD = governance.USD
	g.DAOTreasury = governance.DAOTreasury
	g.Investors = governance.Investors
	g.InitialLidoDevelopers = governance.InitialLidoDevelopers
	g.Foundation = governance.Foundation
	g.Validators = governance.Validators
	return g
}
