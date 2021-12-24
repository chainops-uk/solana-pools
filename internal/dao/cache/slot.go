package cache

import (
	"fmt"
)

const slotArrKey = "slot_arr_key"

func (c *Cache) SetSlotArr(slots [24]float64) {
	c.cache.Set(slotArrKey, slots, storageTime)
}

func (c *Cache) GetSlotArr() ([24]float64, error) {
	v, b := c.cache.Get(slotArrKey)
	if !b {
		return [24]float64{}, fmt.Errorf("%w: %s", KeyWasNotFound, slotArrKey)
	}

	return v.([24]float64), nil
}
