package commands

import (
	"constellation/internal/cli/status"
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(statusCmd)

}

func executeStatus(cmd *cobra.Command, args []string) {

	url := "https://l0-lb-mainnet.constellationnetwork.io/"

	if len(args) > 1 {
		url = args[1]
	}

	status.NewStatus(url).ProvideAsciiStatus(args[0])
}

var statusCmd = &cobra.Command{
	Use:   "status [node-id] ([lb-url])",
	Short: "Node status in ASCII",
	Args:  cobra.RangeArgs(1, 2), // replace with url validation
	Run: func(cmd *cobra.Command, args []string) {
		executeStatus(cmd, args)
	},
}
