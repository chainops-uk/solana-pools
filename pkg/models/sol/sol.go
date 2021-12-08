package sol

import "github.com/shopspring/decimal"

type SOL struct {
	decimal.Decimal
}

func (sol *SOL) SetLamports(lamports uint64) *SOL {
	sol.Decimal = decimal.New(int64(lamports), -9)
	return sol
}

func (sol *SOL) ToLamports() uint64 {
	return uint64(sol.IntPart())
}
