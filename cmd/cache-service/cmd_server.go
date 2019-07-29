package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"

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

			store := store.NewFromEnvironment(log)

			port := env.FromEnvWithDefaultStr("PORT", "3000")
			binding := fmt.Sprintf(":%s", port)

			go configureHealthChecks(log, store)

			serv := service.NewCacheService(log, store)

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

// healthCheckable constrains the type accepted to the registerHealthChecks to
// instances that implement a Readiness and a Liveness probe
type healthCheckable interface {
	ReadinessCheck() (string, func() error)
	LivenessCheck() (string, func() error)
}

func configureHealthChecks(logger zerolog.Logger, healthCheckables ...healthCheckable) {

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

	for _, h := range healthCheckables {

		key, f := h.LivenessCheck()
		health.AddLivenessCheck(key, f)

		key, f = h.ReadinessCheck()
		health.AddReadinessCheck(key, f)
	}

	// nolint : errcheck
	go http.ListenAndServe(":8086", health)

}
