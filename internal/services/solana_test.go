package services_test

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/cache"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"github.com/shopspring/decimal"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestGetPrice(t *testing.T) {
	data := map[string]struct {
		DAO    services.Imp
		Result decimal.Decimal
		Err    error
	}{
		"first": {
			Result: decimal.Decimal{},
			Err:    fmt.Errorf("%w: %s", cache.KeyWasNotFound, "price_key"),
		},
	}

	for s, s2 := range data {
		s2.DAO.Cache = cache.New(time.Minute, time.Minute)
		t.Run(s, func(t *testing.T) {
			price, err := s2.DAO.GetPrice()
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			t.Run(fmt.Sprintf("price[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, price, s2.Result)
			})

		})
	}
}

func TestGetAPY(t *testing.T) {
	data := map[string]struct {
		DAO    services.Imp
		Result decimal.Decimal
		Err    error
	}{
		"first": {
			Result: decimal.Decimal{},
			Err:    fmt.Errorf("%w: %s", cache.KeyWasNotFound, "apy_key"),
		},
	}

	for s, s2 := range data {
		s2.DAO.Cache = cache.New(time.Minute, time.Minute)
		t.Run(s, func(t *testing.T) {
			apy, err := s2.DAO.GetAPY()
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			t.Run(fmt.Sprintf("apy[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, apy, s2.Result)
			})

		})
	}
}

func TestGetValidators(t *testing.T) {
	data := map[string]struct {
		DAO    services.Imp
		Result int64
		Err    error
	}{
		"first": {
			Result: 0,
			Err:    fmt.Errorf("%w: %s", cache.KeyWasNotFound, "validators_key"),
		},
	}

	for s, s2 := range data {
		s2.DAO.Cache = cache.New(time.Minute, time.Minute)
		t.Run(s, func(t *testing.T) {
			val, err := s2.DAO.GetValidators()
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			t.Run(fmt.Sprintf("validator[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, val, s2.Result)
			})

		})
	}
}

func TestGetActiveStake(t *testing.T) {
	data := map[string]struct {
		DAO    services.Imp
		Result uint64
	}{
		"first": {
			Result: 0,
		},
	}

	for s, s2 := range data {
		s2.DAO.Cache = cache.New(time.Minute, time.Minute)
		t.Run(s, func(t *testing.T) {
			as := s2.DAO.GetActiveStake()
			t.Run(fmt.Sprintf("active stake[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, as, s2.Result)
			})

		})
	}
}

func TestGetEpoch(t *testing.T) {
	data := map[string]struct {
		DAO    services.Imp
		Result *smodels.EpochInfo
		Err    error
	}{
		"first": {
			Result: nil,
			Err:    fmt.Errorf("%w: %s", cache.KeyWasNotFound, "epoch_key"),
		},
	}

	for s, s2 := range data {
		s2.DAO.Cache = cache.New(time.Minute, time.Minute)
		t.Run(s, func(t *testing.T) {
			epoch, err := s2.DAO.GetEpoch()
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			t.Run(fmt.Sprintf("epoch[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, epoch, s2.Result)
			})

		})
	}
}
