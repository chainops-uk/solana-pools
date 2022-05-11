package services

import (
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/shopspring/decimal"
)

func (s Imp) GetPrice() (decimal.Decimal, error) {
	price, err := s.Cache.GetPrice()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return price, nil
}

func (s Imp) GetAPY() (decimal.Decimal, error) {
	apy, err := s.Cache.GetAPY()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return apy, nil
}

func (s Imp) GetValidators() (int64, error) {
	c, err := s.Cache.GetValidatorsCount()
	if err != nil {
		return 0, err
	}

	return c, nil
}

func (s Imp) GetActiveStake() uint64 {
	c, err := s.Cache.GetActiveStake()
	if err != nil {
		return 0
	}

	return c
}

func (s Imp) GetEpoch() (*smodels.EpochInfo, error) {
	c, err := s.Cache.GetCurrentEpochInfo()
	if err != nil {
		return nil, err
	}

	return c, nil
}
