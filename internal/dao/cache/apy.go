package cache

import (
	"fmt"
	"github.com/shopspring/decimal"
)

const apyKey = "apy_key"

func (c *Cache) SetAPY(apy decimal.Decimal) {
	c.cache.Set(apyKey, apy, storageTime)
}

func (c *Cache) GetAPY() (decimal.Decimal, error) {
	v, b := c.cache.Get(apyKey)
	if !b {
		return decimal.Decimal{}, fmt.Errorf("%w: %s", keyWasNotFound, apyKey)
	}

	return v.(decimal.Decimal), nil
}
