package commands

import (
	"constellation/internal/updater"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(executeSelfUpgradeCmd)
	executeSelfUpgradeCmd.Flags().StringP("version", "v", "latest", "target version")
	executeSelfUpgradeCmd.Flags().BoolP("force", "f", false, "force upgrade")
}

func executeSelfUpgrade(cmd *cobra.Command, args []string) {
	version, _ := cmd.Flags().GetString("version")
	// TODO: if version is latest then handle it

	// TODO: check if current version is matching our target; abandon if not forced

	updater.SelfUpgrade(version).Run()
}

var executeSelfUpgradeCmd = &cobra.Command{
	Use:   "self-upgrade",
	Short: "Selfupgrade of upgrade manager",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}