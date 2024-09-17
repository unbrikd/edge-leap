package elcli

import (
	"github.com/spf13/cobra"
)

var draftCmd = &cobra.Command{
	Use:   "draft",
	Short: "Draft allows you to handle a new module deployment locally",
	Run: func(cmd *cobra.Command, args []string) {
		executeDraft()
	},
}

func init() {
	rootCmd.AddCommand(draftCmd)
}

func executeDraft() {
	rootCmd.Help()
}
