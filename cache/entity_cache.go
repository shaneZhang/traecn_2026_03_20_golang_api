package cache

import (
	"time"

	"github.com/casdoor/casdoor/object"
)

type UserCache interface {
	Get(id string) (*object.User, bool)
	Set(id string, user *object.User, ttl time.Duration)
	Delete(id string)
	Clear()
}

type userCache struct {
	cache *MemoryCache
}

func NewUserCache() UserCache {
	return &userCache{
		cache: GetMemoryCache(),
	}
}

func (c *userCache) Get(id string) (*object.User, bool) {
	val, found := c.cache.Get(GetUserCacheKey(id))
	if !found {
		return nil, false
	}
	user, ok := val.(*object.User)
	return user, ok
}

func (c *userCache) Set(id string, user *object.User, ttl time.Duration) {
	if ttl == 0 {
		ttl = UserCacheTTL
	}
	c.cache.Set(GetUserCacheKey(id), user, ttl)
}

func (c *userCache) Delete(id string) {
	c.cache.Delete(GetUserCacheKey(id))
}

func (c *userCache) Clear() {
	c.cache.Clear()
}

type OrganizationCache interface {
	Get(id string) (*object.Organization, bool)
	Set(id string, org *object.Organization, ttl time.Duration)
	Delete(id string)
	Clear()
}

type organizationCache struct {
	cache *MemoryCache
}

func NewOrganizationCache() OrganizationCache {
	return &organizationCache{
		cache: GetMemoryCache(),
	}
}

func (c *organizationCache) Get(id string) (*object.Organization, bool) {
	val, found := c.cache.Get(GetOrganizationCacheKey(id))
	if !found {
		return nil, false
	}
	org, ok := val.(*object.Organization)
	return org, ok
}

func (c *organizationCache) Set(id string, org *object.Organization, ttl time.Duration) {
	if ttl == 0 {
		ttl = OrgCacheTTL
	}
	c.cache.Set(GetOrganizationCacheKey(id), org, ttl)
}

func (c *organizationCache) Delete(id string) {
	c.cache.Delete(GetOrganizationCacheKey(id))
}

func (c *organizationCache) Clear() {
	c.cache.Clear()
}

type ApplicationCache interface {
	Get(id string) (*object.Application, bool)
	GetByClientId(clientId string) (*object.Application, bool)
	Set(id string, app *object.Application, ttl time.Duration)
	SetByClientId(clientId string, app *object.Application, ttl time.Duration)
	Delete(id string)
	DeleteByClientId(clientId string)
	Clear()
}

type applicationCache struct {
	cache *MemoryCache
}

func NewApplicationCache() ApplicationCache {
	return &applicationCache{
		cache: GetMemoryCache(),
	}
}

func (c *applicationCache) Get(id string) (*object.Application, bool) {
	val, found := c.cache.Get(GetApplicationCacheKey(id))
	if !found {
		return nil, false
	}
	app, ok := val.(*object.Application)
	return app, ok
}

func (c *applicationCache) GetByClientId(clientId string) (*object.Application, bool) {
	val, found := c.cache.Get(GetApplicationCacheKey("client:" + clientId))
	if !found {
		return nil, false
	}
	app, ok := val.(*object.Application)
	return app, ok
}

func (c *applicationCache) Set(id string, app *object.Application, ttl time.Duration) {
	if ttl == 0 {
		ttl = AppCacheTTL
	}
	c.cache.Set(GetApplicationCacheKey(id), app, ttl)
}

func (c *applicationCache) SetByClientId(clientId string, app *object.Application, ttl time.Duration) {
	if ttl == 0 {
		ttl = AppCacheTTL
	}
	c.cache.Set(GetApplicationCacheKey("client:"+clientId), app, ttl)
}

func (c *applicationCache) Delete(id string) {
	c.cache.Delete(GetApplicationCacheKey(id))
}

func (c *applicationCache) DeleteByClientId(clientId string) {
	c.cache.Delete(GetApplicationCacheKey("client:" + clientId))
}

func (c *applicationCache) Clear() {
	c.cache.Clear()
}
