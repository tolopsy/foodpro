package provider

import (
	"github.com/tolopsy/foodpro/api/persistence"
	"github.com/tolopsy/foodpro/api/persistence/cache"
	"github.com/tolopsy/foodpro/api/persistence/cache/redisclient"
)

type CACHE_SERVER string

const (
	REDIS CACHE_SERVER = "redis"
)

func NewCacheHandler(cacheType, host, password string) (persistence.CacheHandler, error) {
	switch CACHE_SERVER(cacheType) {
	case REDIS:
		return redisclient.NewCacheHandler(host, password)
	default:
		return nil, cache.ErrorCacheServerPluginDoesNotExist
	}
}