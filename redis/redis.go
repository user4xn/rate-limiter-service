package redis

import (
	"rate-limiter/pkg/util"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func GetRedisClient() *redis.Client {
	redisDb, _ := strconv.Atoi(util.GetEnv("REDIS_DB", "1"))
	client := redis.NewClient(&redis.Options{
		Addr:     util.GetEnv("REDIS_URL", "localhost:6379"),
		Password: util.GetEnv("REDIS_PASS", ""),
		DB:       redisDb,
	})

	return client
}
