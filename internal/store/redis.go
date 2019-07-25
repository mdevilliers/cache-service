package store

import (
	"time"

	"github.com/go-redis/redis"
)

type redisstore struct {
	client *redis.Client
}

// NewRedisStore wraps a configured redis client
func NewRedisStore(client *redis.Client) *redisstore {
	return &redisstore{
		client: client,
	}
}

func (r *redisstore) Set(key string, contents string, ttl time.Duration) error {
	return r.client.Set(key, contents, ttl).Err()

}
func (r *redisstore) Get(key string) (string, error) {
	return r.client.Get(key).Result()

}
func (r *redisstore) RandomKey() (string, error) {
	return r.client.RandomKey().Result()
}
func (r *redisstore) Del(key string) error {
	return r.client.Del(key).Err()
}
