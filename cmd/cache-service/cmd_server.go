package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/heptiolabs/healthcheck"
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

			for _, e := range os.Environ() {
				log.Error().Msg(e)
			}

			version.Write(os.Stdout)

			redisServer, _ := os.LookupEnv("REDIS_MASTER_SERVICE_HOST")
			redisPort, _ := os.LookupEnv("REDIS_MASTER_SERVICE_PORT")

			redisClient := redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%s", redisServer, redisPort),
			})

			port, _ := os.LookupEnv("PORT")
			binding := fmt.Sprintf(":%s", port)

			go configureHealthChecks(log, redisClient)

			serv := service.NewCacheService(log, store.NewRedisStore(redisClient))

			server := server.New(log, serv)
			err := server.Start(binding)
			if err != nil {
				return err
			}

			// TODO : add graceful shutdown
			defer server.Stop()

			return nil
		},
	}
	root.AddCommand(cmd)
}

func configureHealthChecks(logger zerolog.Logger, redisClient *redis.Client) {

	health := healthcheck.NewHandler()

	health.AddLivenessCheck("random-failure", func() error {

		r := rand.Intn(100)

		logger.Info().Fields(map[string]interface{}{
			"random": r,
		}).Msg("random failure called")
		if r == 1 {
			return errors.New("boom")
		}
		return nil

	})

	health.AddReadinessCheck("redis-service-check", func() error {
		_, err := redisClient.Ping().Result()
		return err
	})

	go http.ListenAndServe(":8086", health)

}
