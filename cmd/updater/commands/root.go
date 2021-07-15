package commands

import (
"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "constellation-commands-update",
		Short: "Constellation Command Line Utility Update Manager",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
}