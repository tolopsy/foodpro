package redisclient

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/tolopsy/foodpro/persistence"
	"github.com/tolopsy/foodpro/persistence/cache"
)

type CacheHandler struct {
	client    *redis.Client
	recipeKey string
}

func NewCacheHandler(host, password string) (*CacheHandler, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		return nil, err
	}

	return &CacheHandler{client: redisClient, recipeKey: "recipes"}, nil
}

func (handler *CacheHandler) SetRecipes(recipes []persistence.Recipe) error {
	data, err := json.Marshal(recipes)
	if err != nil {
		return err
	}
	return handler.client.Set(handler.recipeKey, string(data), 0).Err()
}

func (handler *CacheHandler) GetRecipes() ([]persistence.Recipe, error) {
	value, err := handler.client.Get(handler.recipeKey).Result()

	if err == redis.Nil {
		return nil, cache.ErrorKeyDoesNotExist
	} else if err != nil {
		return nil, err
	}

	recipes := make([]persistence.Recipe, 0)
	err = json.Unmarshal([]byte(value), &recipes)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (handler *CacheHandler) ClearRecipes() error {
	return handler.client.Del("recipes").Err()
}