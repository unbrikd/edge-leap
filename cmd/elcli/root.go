package elcli

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unbrikd/edge-leap/internal/configuration"
)

const DEFAULT_CONFIG_FILE = "./edge-leap.yaml"

var cfgFile string
var force bool
var config configuration.Configuration
var envFlag []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "elcli",
	Short: "The edge leap (el) cli is a tool to streamline the development of edge computing applications.",
	Long: `The edge leap client (elcli) is a tool to streamline the development of edge computing applications.
unbrikd (c) 2024`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by elcli.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", DEFAULT_CONFIG_FILE, "configuration file")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "force an action")
}

func loadConfig() (*configuration.Configuration, error) {
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
