package cmd

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

		executeNewDraft()
	},
}

func init() {
	rootCmd.AddCommand(draftCmd)
	draftCmd.AddCommand(newDraftCmd)

	// Module configuration
	newDraftCmd.Flags().StringVar(&config.Deployment.Id, "id", "", "development deployment ID")
	newDraftCmd.MarkFlagRequired("id")
	viper.BindPFlag("deployment.id", newDraftCmd.Flags().Lookup("id"))
}

func executeDraft() {
	rootCmd.Help()
}

func executeNewDraft() {
	id := strings.Split(uuid.New().String(), "-")[4]
	viper.Set("session", id)
	viper.Set("version", configuration.CONFIG_VERSION)
	viper.Set("deployment.id", config.Deployment.Id)
	viper.WriteConfig()

	fmt.Printf("New draft session was initialized: %s\n", id)
	fmt.Printf("Please edit the file '%s' accordingly.\n", cfgFile)
}
