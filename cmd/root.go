package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unbrikd/edge-leap/internal/configuration"
	"github.com/unbrikd/edge-leap/internal/utils"
)

const DEFAULT_CONFIG_FILE = "./edge-leap.yaml"

var cfgFile string
var force bool
var config configuration.Configuration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elcli",
	Short: "The edge leap (el) cli is a tool to streamline the development of edge computing applications.",
	Long: `The edge leap client (elcli) is a tool to streamline the development of edge computing applications.
unbrikd (c) 2024`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", utils.GetEnv("EL_CONFIG_FILE", DEFAULT_CONFIG_FILE), "configuration file")

	rootCmd.PersistentFlags().BoolVar(&force, "force", false, "force the command to proceed")
}

func initConfig() {
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n", cfgFile)
	} else {
		fmt.Println("No configuration file found, using flags only")
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

}
