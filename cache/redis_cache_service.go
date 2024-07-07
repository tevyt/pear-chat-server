package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheService struct {
	redisClient *redis.Client
	ctx         *context.Context
}

func NewRedisCacheService(redisClient *redis.Client, ctx *context.Context) *RedisCacheService {
	return &RedisCacheService{
		redisClient: redisClient,
		ctx:         ctx,
	}
}

func (cache *RedisCacheService) Put(key string, value string) error {
	redisTTL := time.Duration(30) * time.Minute
	status := cache.redisClient.Set(*cache.ctx, key, value, redisTTL)

	return status.Err()
}

func (cache *RedisCacheService) Get(key string) (string, error) {
	status := cache.redisClient.Get(*cache.ctx, key)

	return status.Val(), status.Err()
}
