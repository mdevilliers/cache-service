package store

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/mdevilliers/cache-service/internal/env"
	"github.com/rs/zerolog"
)

type redisstore struct {
	rwClient *redis.Client
	rClient  *redis.Client
	logger   zerolog.Logger
}

// NewFromEnvironmet creates a store using K8s overrides to locate a
// Read / Write client for Set, Get and Del
// Read client for RandomKey operations
func NewFromEnvironment(logger zerolog.Logger) *redisstore {

	masterRedisServer := env.FromEnvWithDefaultStr("REDIS_MASTER_SERVICE_HOST", "0.0.0.0")
	masterRedisPort := env.FromEnvWithDefaultStr("REDIS_MASTER_SERVICE_PORT", "6379")

	// TODO : get from file
	// garden.io doesn't support secrets as volumns
	password := env.FromEnvWithDefaultStr("REDIS_PASSWORD", "")

	masterRedisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", masterRedisServer, masterRedisPort),
		Password: password,
	})

	slaveRedisServer := env.FromEnvWithDefaultStr("REDIS_SLAVE_SERVICE_HOST", "0.0.0.0")
	slaveRedisPort := env.FromEnvWithDefaultStr("REDIS_SLAVE_SERVICE_PORT", "6379")

	slaveRedisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", slaveRedisServer, slaveRedisPort),
		Password: password,
	})

	return &redisstore{
		rwClient: masterRedisClient,
		rClient:  slaveRedisClient,
		logger:   logger,
	}
}

func (r *redisstore) Set(key string, contents string, ttl time.Duration) error {
	return r.rwClient.Set(key, contents, ttl).Err()
}

func (r *redisstore) Get(key string) (string, error) {
	return r.rwClient.Get(key).Result()
}

func (r *redisstore) RandomKey() (string, error) {
	return r.rClient.RandomKey().Result()
}

func (r *redisstore) Del(key string) error {
	return r.rwClient.Del(key).Err()
}

func (r *redisstore) ReadinessCheck() (string, func() error) {
	return r.doCheck()
}

func (r *redisstore) LivenessCheck() (string, func() error) {
	return r.doCheck()
}

// For the redis store to be operational we need to be able to
// contact the master and slave
func (r *redisstore) doCheck() (string, func() error) {

	return "redis-service-check", func() error {

		_, err := r.rwClient.Ping().Result()

		if err != nil {
			r.logger.Err(err).Msg("failed to contact redis master")
			return err
		}

		_, err = r.rClient.Ping().Result()

		if err != nil {
			r.logger.Err(err).Msg("failed to contact redis slave")
			return err
		}

		return nil
	}
}
