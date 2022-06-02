package services_test

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/dao/cache"
	"github.com/everstake/solana-pools/internal/services"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestGetAvgSlotTimeMS(t *testing.T) {
	data := map[string]struct {
		DAO    services.Imp
		Result float64
		Err    error
	}{
		"first": {
			Result: 550,
			Err:    nil,
		},
	}

	for s, s2 := range data {
		s2.DAO.Cache = cache.New(time.Minute, time.Minute)
		t.Run(s, func(t *testing.T) {
			avg, err := s2.DAO.GetAvgSlotTimeMS()
			if err != nil {
				assert.Equal(t, err.Error(), s2.Err.Error())
				return
			}
			t.Run(fmt.Sprintf("statistics[%s]", s), func(t *testing.T) {
				assert.DeepEqual(t, avg, s2.Result)
			})

		})
	}
}
