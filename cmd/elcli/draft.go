package elcli

import (
	"github.com/spf13/cobra"
)

var draftCmd = &cobra.Command{
	Use:   "draft",
	Short: "Start a new module draft session",
	Run: func(cmd *cobra.Command, args []string) {
		executeDraft()
	},
}

func init() {
	draftCmd.AddCommand(draftDeployCmd)
	draftCmd.AddCommand(draftNewCmd)
}

func executeDraft() {
	rootCmd.Help()
}
