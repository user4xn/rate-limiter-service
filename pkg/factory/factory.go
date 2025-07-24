package factory

import (
	"rate-limiter/pkg/repository"
	"rate-limiter/redis"
)

type Factory struct {
	RedisRepository repository.RedisRepositoryInterface
}

func NewFactory() *Factory {
	rdb := redis.GetRedisClient()
	return &Factory{
		RedisRepository: repository.NewRedisRepository(rdb),
	}
}
