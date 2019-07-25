package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-redis/redis"
	"github.com/heptiolabs/healthcheck"
	"github.com/mdevilliers/cache-service/internal/env"
	"github.com/mdevilliers/cache-service/internal/server"
	"github.com/mdevilliers/cache-service/internal/service"
	"github.com/mdevilliers/cache-service/internal/store"
	"github.com/mdevilliers/cache-service/internal/version"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func registerServerCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Runs the service as a GRPC service",
		RunE: func(cmd *cobra.Command, args []string) error {

			version.Write(os.Stdout)

			redisServer := env.FromEnvWithDefaultStr("REDIS_MASTER_SERVICE_HOST", "0.0.0.0")
			redisPort := env.FromEnvWithDefaultStr("REDIS_MASTER_SERVICE_PORT", "6379")

			redisClient := redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%s", redisServer, redisPort),
			})

			port := env.FromEnvWithDefaultStr("PORT", "3000")
			binding := fmt.Sprintf(":%s", port)

			go configureHealthChecks(log, redisClient)

			serv := service.NewCacheService(log, store.NewRedisStore(redisClient))

			server := server.New(log, serv)

			go func() {
				err := server.Start(binding)
				if err != nil {
					log.Err(err).Msg("error running server")
				}
			}()

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)

			<-stop

			server.Stop()

			return nil
		},
	}
	root.AddCommand(cmd)
}

func configureHealthChecks(logger zerolog.Logger, redisClient *redis.Client) {

	health := healthcheck.NewHandler()

	// life is too easy without some random failures
	// This is only here for demonstration purposes
	health.AddLivenessCheck("random-failure", func() error {

		r := rand.Intn(1000)

		logger.Info().Fields(map[string]interface{}{
			"random": r,
		}).Msg("random failure called")
		if r == 1 {
			return errors.New("boom")
		}
		return nil

	})

	// cache-service is only ready if it can reach redis
	health.AddReadinessCheck("redis-service-check", func() error {
		_, err := redisClient.Ping().Result()

		if err != nil {
			logger.Err(err).Msg("failed to ping redis")
		}

		return err
	})

	// nolint : errcheck
	go http.ListenAndServe(":8086", health)

}
