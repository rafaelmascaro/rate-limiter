package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	Client *redis.Client
}

func NewRedisRepository(
	redisHost string,
	redisPort string,
	db int,
) *RedisRepository {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   db,
	})

	return &RedisRepository{
		Client: redisClient,
	}
}

func (r *RedisRepository) AddKey(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Set(ctx, key, 1, expiration).Err()
}

func (r *RedisRepository) Exists(ctx context.Context, key string) (int64, error) {
	return r.Client.Exists(ctx, key).Result()
}

func (r *RedisRepository) Increment(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

func (r *RedisRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}

func (r *RedisRepository) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisRepository) AddHash(
	ctx context.Context,
	key string,
	limit int,
	timeBlock time.Duration,
) error {
	exists, _ := r.Exists(ctx, key)
	if exists > 0 {
		return nil
	}
	return r.Client.HSet(ctx, key, map[string]interface{}{
		"limit":      limit,
		"time_block": int(timeBlock.Seconds()),
	}).Err()
}

func (r *RedisRepository) Find(ctx context.Context, key string) (bool, int, time.Duration, error) {
	limitStr, err := r.Client.HGet(ctx, key, "limit").Result()
	if err == redis.Nil {
		return false, 0, 0, nil
	} else if err != nil {
		return false, 0, 0, err
	}

	timeBlockStr, err := r.Client.HGet(ctx, key, "time_block").Result()
	if err != nil {
		return false, 0, 0, err
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return false, 0, 0, err
	}

	timeBlock, err := strconv.Atoi(timeBlockStr)
	if err != nil {
		return false, 0, 0, err
	}

	return true, limit, time.Duration(timeBlock) * time.Second, nil
}
