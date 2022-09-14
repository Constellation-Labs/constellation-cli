package commands

import (
	nodemap "constellation/internal/cli/nodemap"
	"constellation/pkg/node"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	home := strings.TrimRight(os.Getenv("HOME"), "/")

	rootCmd.AddCommand(nodeMapCmd)
	nodeMapCmd.Flags().BoolP("silent", "s", false, "run in silent mode")
	nodeMapCmd.Flags().StringP("image", "i", "", "image file path for graphical output")
	nodeMapCmd.Flags().StringP("theme", "t", "transparent", "background theme for image output [light/dark]")
	nodeMapCmd.Flags().BoolP("verbose", "v", false, "provide more detailed output")
	nodeMapCmd.Flags().StringP("operators", "o", fmt.Sprintf("%s/operators", home), "operators file in csv format")
}

func executeNodemap(cmd *cobra.Command, args []string) {
	addr := args[0]
	silent, _ := cmd.Flags().GetBool("silent")

	ng := nodemap.NewNodemap()

	ng.DiscoverNetwork(node.AddrOf(addr), silent)
}

var nodeMapCmd = &cobra.Command{
	Use:   "nodemap [url]",
	Short: "Perform recursive discovery and build network status",
	Args:  cobra.ExactArgs(1), // replace with url validation
	Run: func(cmd *cobra.Command, args []string) {
		executeNodemap(cmd, args)
	},
}
