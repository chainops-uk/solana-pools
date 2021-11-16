package cache

import (
	"fmt"
	"github.com/shopspring/decimal"
)

const priceKey = "price_key"

func (c *Cache) SetPrice(apy decimal.Decimal) {
	c.cache.Set(priceKey, apy, storageTime)
}

func (c *Cache) GetPrice() (decimal.Decimal, error) {
	v, b := c.cache.Get(priceKey)
	if !b {
		return decimal.Decimal{}, fmt.Errorf("%w: %s", keyWasNotFound, priceKey)
	}

	return v.(decimal.Decimal), nil
}
