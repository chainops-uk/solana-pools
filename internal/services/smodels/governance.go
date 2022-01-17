package smodels

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
)

type Governance struct {
	Name               string
	Symbol             string
	VoiceURL           string
	WebSiteURL         string
	Image              string
	GeckoKey           string
	Blockchain         string
	ContractAddress    string
	MaximumTokenSupply float64
	CirculatingSupply  float64
	USD                float64
}

func (g *Governance) Set(governance *dmodels.Governance) *Governance {
	g.Name = governance.Name
	g.Symbol = governance.Symbol
	g.VoiceURL = governance.VoiceURL
	g.WebSiteURL = governance.WebSiteURL
	g.Image = governance.Image
	g.GeckoKey = governance.GeckoKey
	g.Blockchain = governance.Blockchain
	g.ContractAddress = governance.ContractAddress
	g.MaximumTokenSupply = governance.MaximumTokenSupply
	g.CirculatingSupply = governance.CirculatingSupply
	g.USD = governance.USD
	return g
}
