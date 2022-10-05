package commands

import (
	nodegrid2 "constellation/internal/cli/nodegrid"
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(nodegridCmd)
	nodegridCmd.Flags().BoolP("silent", "s", false, "run in silent mode")
	nodegridCmd.Flags().StringP("image", "i", "", "image file path for graphical output")
	nodegridCmd.Flags().StringP("theme", "t", "transparent", "background theme for image output [light/dark]")
	nodegridCmd.Flags().BoolP("verbose", "v", false, "provide more detailed output")
}

func executeNodegrid(cmd *cobra.Command, args []string) {

	url := "https://l0-lb-mainnet.constellationnetwork.io/"

	if len(args) > 0 {
		url = args[0]
	}

	silent, _ := cmd.Flags().GetBool("silent")
	outputImage, _ := cmd.Flags().GetString("image")
	outputTheme, _ := cmd.Flags().GetString("theme")
	verbose, _ := cmd.Flags().GetBool("verbose")
	operatorsFile, _ := cmd.Flags().GetString("operators")

	ng := nodegrid2.NewNodegrid(operatorsFile)

	ng.BuildNetworkStatus(url, silent, outputImage, outputTheme, verbose)
}

var nodegridCmd = &cobra.Command{
	Use:   "nodegrid [url]",
	Short: "Build and verify Constellation Hypergraph Network status for a given loadbalancer status url. If not provided mainnet 2.0 lb is used",
	Args:  cobra.RangeArgs(1, 2), // replace with url validation
	Run: func(cmd *cobra.Command, args []string) {
		executeNodegrid(cmd, args)
	},
}
