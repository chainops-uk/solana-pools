package smodels

import "github.com/everstake/solana-pools/internal/dao/dmodels"

type Coin struct {
	Name       string
	Address    string
	USD        float64
	ThumbImage string
	SmallImage string
	LargeImage string
	DeFi       []*DeFi
}

func (c *Coin) Set(coin *dmodels.Coin, fi []*DeFi) *Coin {
	if fi != nil {
		c.DeFi = fi
	}
	c.USD = coin.USD
	c.ThumbImage = coin.ThumbImage
	c.SmallImage = coin.SmallImage
	c.LargeImage = coin.LargeImage
	c.Name = coin.Name
	c.Address = coin.Address
	return c
}
