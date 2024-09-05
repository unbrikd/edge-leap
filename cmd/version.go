package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unbrikd/edge-leap/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of the edge leap client",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s, build %s\n", version.Version, version.Revision)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
