package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	home := strings.TrimRight(os.Getenv("HOME"), "/")

	rootCmd.AddCommand(executeUpgradeCmd)
	executeUpgradeCmd.Flags().BoolP("silent", "s", false, "run in silent mode")
	executeUpgradeCmd.Flags().StringP("image", "i", "", "image file path for graphical output")
	executeUpgradeCmd.Flags().StringP("theme", "t", "transparent", "background theme for image output [light/dark]")
	executeUpgradeCmd.Flags().BoolP("verbose", "v", false, "provide more detailed output")
	executeUpgradeCmd.Flags().StringP("operators", "o", fmt.Sprintf("%s/operators", home), "operators file in csv format")
}

func executeUpgrade(cmd *cobra.Command, args []string) {

}

var executeUpgradeCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for upgrade for your Constellation Labs commands tools",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}