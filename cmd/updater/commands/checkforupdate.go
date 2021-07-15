package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	home := strings.TrimRight(os.Getenv("HOME"), "/")

	rootCmd.AddCommand(checkForUpgradeCmd)
	checkForUpgradeCmd.Flags().BoolP("silent", "s", false, "run in silent mode")
	checkForUpgradeCmd.Flags().StringP("image", "i", "", "image file path for graphical output")
	checkForUpgradeCmd.Flags().StringP("theme", "t", "transparent", "background theme for image output [light/dark]")
	checkForUpgradeCmd.Flags().BoolP("verbose", "v", false, "provide more detailed output")
	checkForUpgradeCmd.Flags().StringP("operators", "o", fmt.Sprintf("%s/operators", home), "operators file in csv format")
}

func executeCheckForUpdate(cmd *cobra.Command, args []string) {

}

var checkForUpgradeCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for upgrade for your Constellation Labs commands tools",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}