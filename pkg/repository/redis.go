package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepositoryInterface interface {
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, duration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)

	Set(ctx context.Context, key string, value any, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type redisRepository struct {
	Rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) *redisRepository {
	return &redisRepository{
		Rdb: rdb,
	}
}

func (r *redisRepository) Incr(ctx context.Context, key string) (int64, error) {
	val, err := r.Rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r *redisRepository) Expire(ctx context.Context, key string, duration time.Duration) error {
	err := r.Rdb.Expire(ctx, key, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisRepository) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.Rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return ttl, nil
}

func (r *redisRepository) Set(ctx context.Context, key string, value any, duration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.Rdb.Set(ctx, key, jsonData, duration).Err()
	return err
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *redisRepository) Del(ctx context.Context, key string) error {
	err := r.Rdb.Del(ctx, key).Err()
	if err != nil {
		return errors.New("failed delete")
	}
	return nil
}
