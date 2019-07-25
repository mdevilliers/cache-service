package main

import (
	"os"

	"github.com/mdevilliers/cache-service/internal/version"

	"github.com/spf13/cobra"
)

func registerVersionCommand(root *cobra.Command) {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the build version",
		Run: func(*cobra.Command, []string) {
			version.Write(os.Stdout)
		},
	}

	root.AddCommand(cmd)
}
