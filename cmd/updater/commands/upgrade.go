package commands

import (
	"constellation/internal/updater"
	"github.com/spf13/cobra"
	"log"
)

func init() {

	rootCmd.AddCommand(executeUpgradeCmd)
	executeUpgradeCmd.Flags().StringP("version", "v", "latest", "target version")
	executeUpgradeCmd.Flags().BoolP("force", "f", false, "force upgrade")
}

func executeUpgrade(cmd *cobra.Command, args []string) {
	version, _ := cmd.Flags().GetString("version")
	// TODO: if version is latest then handle it

	// TODO: check if current version is matching our target; abandon if not forced

	err := updater.CommandlineUpgrade(version).Run()
	if err != nil {
		log.Fatal(err)
	}
}

var executeUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade for your Constellation Labs commands tools",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}
