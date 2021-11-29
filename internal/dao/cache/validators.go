package cache

import (
	"fmt"
)

const validatorsKey = "validators_key"

func (c *Cache) SetValidatorCount(count int64) {
	c.cache.Set(validatorsKey, count, storageTime)
}

func (c *Cache) GetValidatorsCount() (int64, error) {
	v, b := c.cache.Get(validatorsKey)
	if !b {
		return 0, fmt.Errorf("%w: %s", KeyWasNotFound, validatorsKey)
	}

	return v.(int64), nil
}
