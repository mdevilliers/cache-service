package main

import (
	"io"
	"os"

	"github.com/mdevilliers/cache-service/internal/version"
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
		Short:         "TODO: ???",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Setup default logger
			log = initLogger(logLevel, useConsole, makeVerbose)
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

func initLogger(logLevel string, useConsole, makeVerbose bool) zerolog.Logger {
	// Set logger level field to severity for stack driver support
	zerolog.LevelFieldName = "severity"
	var w io.Writer = os.Stdout
	if useConsole {
		w = zerolog.ConsoleWriter{
			Out: os.Stdout,
		}
	}
	// Parse level from config
	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	// Override level with verbose
	if makeVerbose {
		lvl = zerolog.DebugLevel
	}
	return zerolog.New(w).Level(lvl).With().Fields(map[string]interface{}{
		"version": version.Version,
		"app":     "cache-service",
	}).Timestamp().Logger()
}
