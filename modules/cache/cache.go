package cache

import (
	"time"

	"github.com/dgraph-io/ristretto"
)

type LocalCache struct {
	Cache *ristretto.Cache
}

// Init will initialize the cache, returns an error if error occurs
func (la *LocalCache) Init() error {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	la.Cache = cache
	if err != nil {
		return err
	}
	return nil
}

// Set will put the value into the cache with a TTL of 1 hour
func (la *LocalCache) Set(key string, value interface{}) {
	la.Cache.SetWithTTL(key, value, 1, 1*time.Hour)
}

// Get will return the stored value in the cache by the key provided
func (la *LocalCache) Get(key string) (interface{}, bool) {
	value, wasRetrieved := la.Cache.Get(key)
	return value, wasRetrieved
}
