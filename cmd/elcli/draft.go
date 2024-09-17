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

// draftCmd represents the draft command
var draftCmd = &cobra.Command{
	Use:   "draft",
	Short: "Draft allows you to handle a new module deployment locally",
	Run: func(cmd *cobra.Command, args []string) {
		executeDraft()
	},
}

var newDraftCmd = &cobra.Command{
	Use:   "new",
	Short: "Sets a new draft session for a module",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(cfgFile); err == nil {
			if !force {
				fmt.Printf("Configuration file already exists. Use --force to overwrite\n")
				return
			}

			if err := os.Remove(cfgFile); err != nil {
				fmt.Printf("Error removing configuration file: %v\n", err)
				return
			}

			viper.SetConfigFile(cfgFile)
		}

		initConfig()
		executeNewDraft()
	},
}

func init() {
	rootCmd.AddCommand(draftCmd)
	draftCmd.AddCommand(newDraftCmd)
}

func executeDraft() {
	rootCmd.Help()
}

func executeNewDraft() {
	id := strings.Split(uuid.New().String(), "-")[4]
	viper.Set("session", id)
	viper.Set("version", configuration.CONFIG_VERSION)
	viper.WriteConfig()

	fmt.Println(id)
}
