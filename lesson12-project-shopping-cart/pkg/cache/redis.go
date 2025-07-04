package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheService struct {
	ctx context.Context
	rdb *redis.Client
}

func NewRedisCacheService(rdb *redis.Client) *RedisCacheService {
	return &RedisCacheService{
		ctx: context.Background(),
		rdb: rdb,
	}
}

func (cs *RedisCacheService) Get(key string, dest any) error {
	data, err := cs.rdb.Get(cs.ctx, key).Result()

	if err == redis.Nil {
		return err
	}

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func (cs *RedisCacheService) Set(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return cs.rdb.Set(cs.ctx, key, data, ttl).Err()
}

func (cs *RedisCacheService) Clear(pattern string) error {
	cusor := uint64(0)
	for {
		keys, nextCursor, err := cs.rdb.Scan(cs.ctx, cusor, pattern, 2).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			cs.rdb.Del(cs.ctx, keys...)
		}

		cusor = nextCursor

		if cusor == 0 {
			break
		}
	}

	return nil
}
