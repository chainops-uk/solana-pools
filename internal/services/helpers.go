package services

import "github.com/shopspring/decimal"

const SolanaPrecision = 9

func lampToSol(lamports uint64) decimal.Decimal {
	return decimal.New(int64(lamports), SolanaPrecision)
}
