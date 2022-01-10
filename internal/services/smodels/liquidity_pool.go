package smodels

import "github.com/everstake/solana-pools/internal/dao/dmodels"

type LiquidityPool struct {
	Name  string
	About string
	Image string
	URL   string
}

func (lp *LiquidityPool) Set(pool *dmodels.LiquidityPool) *LiquidityPool {
	lp.Name = pool.Name
	lp.About = pool.About
	lp.URL = pool.URL
	lp.Image = pool.Image
	return lp
}
