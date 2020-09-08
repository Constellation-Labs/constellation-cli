package cmd

import (
	nodegrid "constellation_cli/nodegrid"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	home := strings.TrimRight(os.Getenv("HOME"), "/")

	rootCmd.AddCommand(nodegridCmd)
	nodegridCmd.Flags().BoolP("silent", "s", false, "run in silent mode")
	nodegridCmd.Flags().StringP("image", "i", "", "image file path for graphical output")
	nodegridCmd.Flags().StringP("theme", "t", "transparent", "background theme for image output [light/dark]")
	nodegridCmd.Flags().BoolP("verbose", "v", false, "provide more detailed output")
	nodegridCmd.Flags().StringP("operators", "o", fmt.Sprintf("%s/operators", home), "operators file in csv format")
}

func executeNodegrid(cmd *cobra.Command, args []string) {
	url := args[0]
	silent, _ := cmd.Flags().GetBool("silent")
	outputImage, _ := cmd.Flags().GetString("image")
	outputTheme, _ := cmd.Flags().GetString("theme")
	verbose, _ := cmd.Flags().GetBool("verbose")
	operatorsFile, _ := cmd.Flags().GetString("operators")

	ng := nodegrid.NewNodegrid(operatorsFile)

	ng.BuildNetworkStatus(url, silent, outputImage, outputTheme, verbose)
}

var nodegridCmd = &cobra.Command{
	Use:   "nodegrid [url]",
	Short: "Build and verify Constellation Hypergraph Network status for a given loadbalancer status url",
	Args: cobra.ExactArgs(1), // replace with url validation
	Run: func(cmd *cobra.Command, args []string) {
		executeNodegrid(cmd, args)
	},
}