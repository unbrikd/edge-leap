package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/unbrikd/edge-leap/internal/configuration"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the edge-leap configuration",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(file)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("configuration file does not exist: %s\n", file)
			return
		}

		c, err := configuration.Load(file)
		if err != nil {
			fmt.Printf("failed to load configuration: %v\n", err)
			return
		}

		fmt.Println(c.Id)
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new development session and create a new configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(file)
		if !errors.Is(err, os.ErrNotExist) && !force {
			fmt.Println("configuration file already exists, to overwrite it use the --force flag")
			return
		}

		c, err := configuration.New(file)
		if err != nil {
			fmt.Printf("failed to create configuration file: %v\n", err)
			return
		}

		fmt.Println(c.Id)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
}
