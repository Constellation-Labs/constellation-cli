package commands

import (
	nodemap "constellation/internal/cli/nodemap"
	"constellation/pkg/node"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nodeMapCmd)
	nodeMapCmd.Flags().BoolP("silent", "s", false, "run in silent mode")
	nodeMapCmd.Flags().StringP("image", "i", "", "image file path for graphical output")
	nodeMapCmd.Flags().StringP("theme", "t", "transparent", "background theme for image output [light/dark]")
	nodeMapCmd.Flags().BoolP("verbose", "v", false, "provide more detailed output")
}

func executeNodemap(cmd *cobra.Command, args []string) {
	addr := args[0]
	silent, _ := cmd.Flags().GetBool("silent")
	outputImage, _ := cmd.Flags().GetString("image")
	outputTheme, _ := cmd.Flags().GetString("theme")
	verbose, _ := cmd.Flags().GetBool("verbose")
	ng := nodemap.NewNodemap()

	ng.DiscoverNetwork(node.AddrOf(addr), silent, verbose, outputImage, outputTheme)
}

var nodeMapCmd = &cobra.Command{
	Use:   "nodemap [url]",
	Short: "Perform recursive discovery and build network status",
	Args:  cobra.ExactArgs(1), // replace with url validation
	Run: func(cmd *cobra.Command, args []string) {
		executeNodemap(cmd, args)
	},
}
