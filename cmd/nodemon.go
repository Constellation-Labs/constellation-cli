package cmd

import (
	nodemon "constellation_cli/nodemon"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	home := strings.TrimRight(os.Getenv("HOME"), "/")

	rootCmd.AddCommand(nodemonCmd)

	nodemonCmd.Flags().StringP("operators", "o", fmt.Sprintf("%s/operators", home), "operators file in csv format")

	nodemonCmd.Flags().StringP("configFile", "c", fmt.Sprintf("%s/webhook", home), "webhook url file")

	nodemonCmd.Flags().StringP("statusFile", "s", fmt.Sprintf("%s/network-status", home), "status cache file")

	nodemonCmd.Flags().StringP("theme", "t", "dark", "color theme")
}

func executeNodemon(cmd *cobra.Command, args []string) {
	url := args[0]
	operatorsFile, _ := cmd.Flags().GetString("operators")
	outputTheme, _ := cmd.Flags().GetString("theme")
	configFile, _ := cmd.Flags().GetString("configFile")
	statusFile, _ := cmd.Flags().GetString("statusFile")

	nm := nodemon.NewNodemon()

	nm.ExecuteNodesCheck(url, configFile, statusFile, outputTheme, operatorsFile)
}

var nodemonCmd = &cobra.Command{
	Use:   "nodemon [url]",
	Short: "Build and verify Constellation Hypergraph Network status for a given loadbalancer status url",
	Args: cobra.ExactArgs(1), // replace with url validation
	Run: func(cmd *cobra.Command, args []string) {
		executeNodemon(cmd, args)
	},
}