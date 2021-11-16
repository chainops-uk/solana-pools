package cache

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"time"
)

const storageTime = time.Hour * 24

var (
	keyWasNotFound = errors.New("the key was not found")
)

type Cache struct {
	cache *cache.Cache
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	return &Cache{cache: cache.New(defaultExpiration, cleanupInterval)}
}
