package elcli

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unbrikd/edge-leap/internal/configuration"
)

var newDraftCmd = &cobra.Command{
	Use:   "new",
	Short: "Sets a new draft session for a module",
	Run: func(cmd *cobra.Command, args []string) {
		preExecuteChecksNewDraft()
		LoadConfig()
		executeNewDraft()
	},
}

func init() {
	draftCmd.AddCommand(newDraftCmd)
}

// preExecuteChecksNewDraft checks if the configuration file exists and if the --force flag is set.
// If the configuration file exists and the --force flag is not set, the function will exit with an error message
// If the configuration file exists and the --force flag is set, the function will remove the configuration file and create a new one
func preExecuteChecksNewDraft() {
	if _, err := os.Stat(cfgFile); err == nil {
		if !force {
			fmt.Printf("Configuration file already exists. Use --force to overwrite\n")
			os.Exit(1)
		}

		if err := os.Remove(cfgFile); err != nil {
			fmt.Printf("Error removing configuration file: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Create(cfgFile); err != nil {
			fmt.Printf("Error creating configuration file: %v\n", err)
			os.Exit(1)
		}
	}
}

func executeNewDraft() {
	id := strings.Split(uuid.New().String(), "-")[4]
	viper.Set("session", id)
	viper.Set("version", configuration.CONFIG_VERSION)
	viper.WriteConfig()

	fmt.Println(id)
}
