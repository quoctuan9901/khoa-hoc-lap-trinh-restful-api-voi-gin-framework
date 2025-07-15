package cache

import "time"

type RedisCacheService interface {
	Get(key string, dest any) error
	Set(key string, value any, ttl time.Duration) error
	Clear(pattern string) error
	Exists(key string) (bool, error)
}
