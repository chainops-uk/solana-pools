package cache

import (
	"fmt"
)

const activeStakeKey = "active_stake_key"

func (c *Cache) SetActiveStake(activeStake uint64) {
	c.cache.Set(activeStakeKey, activeStake, storageTime)
}

func (c *Cache) GetActiveStake() (uint64, error) {
	v, b := c.cache.Get(activeStakeKey)
	if !b {
		return 0, fmt.Errorf("%w: %s", keyWasNotFound, activeStakeKey)
	}
	return v.(uint64), nil
}
