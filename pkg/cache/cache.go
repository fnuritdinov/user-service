package cache

import (
	"time"

	memoryCache "github.com/patrickmn/go-cache"
)

type cache struct {
	memory *memoryCache.Cache
}

type CacheMemory struct {
	Name     string
	Email    string
	Password string
	Age      int32
	Phone    string
	Role     string
}

type MemoryCache interface {
	Set(key string, value any, d time.Duration)
	Get(key string) (any, bool)
	Replace(key string, value any, ttl time.Duration) error
	Remove(key string)
}

func NewMemoryCache() MemoryCache {
	c := memoryCache.New(memoryCache.NoExpiration, time.Minute)

	return &cache{memory: c}
}

func (c *cache) Set(key string, value any, ttl time.Duration) {
	c.memory.Set(key, value, ttl)
}

func (c *cache) Get(key string) (any, bool) {
	return c.memory.Get(key)
}

func (c *cache) Replace(key string, value any, ttl time.Duration) error {
	return c.memory.Replace(key, value, ttl)
}

func (c *cache) Remove(key string) {
	c.memory.Delete(key)
}
