package cache

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"time"
)

const totalCurrentStatisticsKey = "total_current_statistics_key"

func (c *Cache) SetCurrentStatistic(statistic *smodels.Statistic, storageTime time.Duration) {
	c.cache.Set(totalCurrentStatisticsKey, statistic, storageTime)
}

func (c *Cache) GetCurrentStatistic() (*smodels.Statistic, error) {
	v, b := c.cache.Get(totalCurrentStatisticsKey)
	if !b {
		return nil, fmt.Errorf("%w: %s", KeyWasNotFound, totalCurrentStatisticsKey)
	}

	return v.(*smodels.Statistic), nil
}
