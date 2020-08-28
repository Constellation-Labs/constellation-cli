package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "constellation_cli",
		Short: "Constellation Command Line Utility",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
}
