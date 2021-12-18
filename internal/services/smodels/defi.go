package smodels

import (
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/shopspring/decimal"
)

type DeFi struct {
	BuyCoin       *Coin
	LiquidityPool *LiquidityPool
	Liquidity     float64
	APY           decimal.Decimal
}

func (f *DeFi) Set(defi *dmodels.DEFI, buyCoin *Coin, liquidityPool *LiquidityPool) *DeFi {
	f.APY = defi.APY
	f.LiquidityPool = liquidityPool
	f.BuyCoin = buyCoin
	f.Liquidity = defi.Liquidity
	return f
}
