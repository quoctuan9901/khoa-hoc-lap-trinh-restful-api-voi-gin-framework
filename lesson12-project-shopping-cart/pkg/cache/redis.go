package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCacheService struct {
	ctx context.Context
	rdb *redis.Client
}

func NewRedisCacheService(rdb *redis.Client) RedisCacheService {
	return &redisCacheService{
		ctx: context.Background(),
		rdb: rdb,
	}
}

func (cs *redisCacheService) Get(key string, dest any) error {
	data, err := cs.rdb.Get(cs.ctx, key).Result()

	if err == redis.Nil {
		return err
	}

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (cs *redisCacheService) Set(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cs.rdb.Set(cs.ctx, key, data, ttl).Err()
}

func (cs *redisCacheService) Clear(pattern string) error {
	cursor := uint64(0)
	for {
		keys, nextCursor, err := cs.rdb.Scan(cs.ctx, cursor, pattern, 2).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			cs.rdb.Del(cs.ctx, keys...)
		}

		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (cs *redisCacheService) Exists(key string) (bool, error) {
	count, err := cs.rdb.Exists(cs.ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
