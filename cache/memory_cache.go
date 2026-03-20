package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type MemoryCache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
}

var (
	instance *MemoryCache
	once     sync.Once
)

func GetMemoryCache() *MemoryCache {
	once.Do(func() {
		instance = &MemoryCache{
			items: make(map[string]*CacheItem),
		}
		go instance.cleanupExpiredItems()
	})
	return instance
}

func (c *MemoryCache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	
	c.items[key] = &CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return nil, false
	}
	
	return item.Value, true
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.items, key)
}

func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items = make(map[string]*CacheItem)
}

func (c *MemoryCache) cleanupExpiredItems() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		for key, item := range c.items {
			if item.Expiration > 0 && now > item.Expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

const (
	CacheKeyUserPrefix         = "user:"
	CacheKeyOrganizationPrefix = "org:"
	CacheKeyApplicationPrefix  = "app:"
	CacheKeyGroupPrefix        = "group:"
	
	DefaultCacheTTL = 5 * time.Minute
	UserCacheTTL    = 10 * time.Minute
	OrgCacheTTL     = 15 * time.Minute
	AppCacheTTL     = 10 * time.Minute
)

func GetUserCacheKey(id string) string {
	return CacheKeyUserPrefix + id
}

func GetOrganizationCacheKey(id string) string {
	return CacheKeyOrganizationPrefix + id
}

func GetApplicationCacheKey(id string) string {
	return CacheKeyApplicationPrefix + id
}

func GetGroupCacheKey(id string) string {
	return CacheKeyGroupPrefix + id
}
