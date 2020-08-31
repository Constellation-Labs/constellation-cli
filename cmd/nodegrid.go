package cmd

import (
	nodegrid "constellation_cli/nodegrid"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nodegridCmd)
	nodegridCmd.Flags().BoolP("silent", "s", false, "run in silent mode")
	nodegridCmd.Flags().StringP("image", "i", "", "image file path for graphical output")
	nodegridCmd.Flags().StringP("theme", "t", "transparent", "background theme for image output [light/dark]")
}

func executeNodegrid(cmd *cobra.Command, args []string) {
	url := args[0]
	silent, _ := cmd.Flags().GetBool("silent")
	outputImage, _ := cmd.Flags().GetString("image")
	outputTheme, _ := cmd.Flags().GetString("theme")

	ng := nodegrid.NewNodegrid()

	ng.BuildNetworkStatus(url, silent, outputImage, outputTheme)
}

var nodegridCmd = &cobra.Command{
	Use:   "nodegrid [url]",
	Short: "Build and verify Constellation Hypergraph Network status for a given loadbalancer status url",
	Args: cobra.ExactArgs(1), // replace with url validation
	Run: func(cmd *cobra.Command, args []string) {
		executeNodegrid(cmd, args)
	},
}