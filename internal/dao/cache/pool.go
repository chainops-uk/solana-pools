package cache

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"time"
)

const PoolKey = "pool_key"

func (c *Cache) SetPool(pool *smodels.PoolDetails, storageTime time.Duration) {
	c.cache.Set(fmt.Sprintf("%s:%s", PoolKey, pool.Pool.Name), pool, storageTime)
}

func (c *Cache) GetPool(name string) (*smodels.PoolDetails, error) {
	v, b := c.cache.Get(fmt.Sprintf("%s:%s", PoolKey, name))
	if !b {
		return nil, fmt.Errorf("%w: %s", KeyWasNotFound, totalCurrentStatisticsKey)
	}

	return v.(*smodels.PoolDetails), nil
}
