package main

import (
	"os"

	"github.com/mdevilliers/cache-service/internal/env"
	"github.com/mdevilliers/cache-service/internal/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// Application entry point
func main() {
	cmd := rootCmd()
	if err := cmd.Execute(); err != nil {
		log.Error().Err(err).Msg("exiting from fatal error")
		os.Exit(1)
	}
}

// Default logger
var log zerolog.Logger

func rootCmd() *cobra.Command {
	useConsole := false
	makeVerbose := false
	logLevel := "INFO"

	cmd := &cobra.Command{
		Use:           "cache-service",
		Short:         "Demonstration GRPC K8s service",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// Setup default logger
			ll := env.FromEnvWithDefaultStr("CACHE_SERVICE_LOG_LEVEL", logLevel)
			uc := env.FromEnvWithDefaultBool("CACHE_SERVICE_LOG_USE_CONSOLE", useConsole)
			mv := env.FromEnvWithDefaultBool("CACHE_SERVICE_LOG_VERBOSE", makeVerbose)

			log = logger.New(ll, uc, mv)
			return nil
		},
	}
	// Global flags
	pflags := cmd.PersistentFlags()
	pflags.BoolVar(&useConsole, "console", useConsole, "use console log writer")
	pflags.BoolVarP(&makeVerbose, "verbose", "v", makeVerbose, "verbose logging")
	pflags.StringVar(&logLevel, "log-level", logLevel, "log level")

	// Add sub commands
	registerVersionCommand(cmd)
	registerServerCommand(cmd)
	registerClientCommand(cmd)

	return cmd
}
