package services

import (
	"github.com/shopspring/decimal"
)

func (s Imp) GetPrice() (decimal.Decimal, error) {
	price, err := s.cache.GetPrice()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return price, nil
}

func (s Imp) GetAPY() (decimal.Decimal, error) {
	apy, err := s.cache.GetAPY()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return apy, nil
}

func (s Imp) GetValidators() (int64, error) {
	c, err := s.cache.GetValidatorsCount()
	if err != nil {
		return 0, err
	}

	return c, nil
}

func (s Imp) GetActiveStake() uint64 {
	c, err := s.cache.GetActiveStake()
	if err != nil {
		return 0
	}

	return c
}
