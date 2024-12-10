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

var draftNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Sets a new module draft configuration",
	PreRun: func(cmd *cobra.Command, args []string) {
		handleConfigFileCreation()
	},
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()
		executeDraftNew()
	},
}

// handleConfigFileCreation checks if the configuration file exists and if the --force flag is set.
// If the configuration file exists and the --force flag is not set, the function will exit with an error message
// If the configuration file exists and the --force flag is set, the function will remove the configuration file and create a new one
func handleConfigFileCreation() {
	if _, err := os.Stat(cfgFile); err == nil {
		if !force {
			fmt.Println("configuration file already exists, use --force to overwrite")
			os.Exit(1)
		}

		if err := os.Remove(cfgFile); err != nil {
			fmt.Printf("error deleting existing configuration file: %v\n", err)
			os.Exit(1)
		}
	}

	if _, err := os.Create(cfgFile); err != nil {
		fmt.Printf("error creating new configuration file: %v\n", err)
		os.Exit(1)
	}
}

// executeNewDraft generates a new draft session by creating a new configuration file to be used to deploy the draft module
func executeDraftNew() {
	id := strings.Split(uuid.New().String(), "-")[4]
	viper.Set("session", id)
	viper.Set("version", configuration.CONFIG_VERSION)
	viper.WriteConfig()

	fmt.Println(id)
}
