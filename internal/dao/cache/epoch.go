package cache

import (
	"fmt"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"time"
)

const epochKey = "epoch_key"

func (c *Cache) GetCurrentEpochInfo() (*smodels.EpochInfo, error) {
	epoch, ok := c.cache.Get(epochKey)
	if !ok {
		return nil, fmt.Errorf("%w: %s", KeyWasNotFound, epochKey)
	}

	return epoch.(*smodels.EpochInfo), nil
}

func (c *Cache) SetCurrentEpochInfo(epoch *smodels.EpochInfo) {
	c.cache.Set(epochKey, epoch, time.Hour*24)
}
