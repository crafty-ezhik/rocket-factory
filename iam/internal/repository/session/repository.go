package session

import "github.com/crafty-ezhik/rocket-factory/platform/pkg/cache"

type repository struct {
	redis cache.RedisClient
}

func NewRepository(redis cache.RedisClient) *repository {
	return &repository{redis: redis}
}
