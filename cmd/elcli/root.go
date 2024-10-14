package elcli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unbrikd/edge-leap/internal/configuration"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const DEFAULT_CONFIG_FILE = "./edge-leap.yaml"

var cfgFile string
var force bool
var config configuration.Configuration
var envFlag []string
var logLevel string

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
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "v", "info", "log level")

	rootCmd.PersistentPreRun = setLogger
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

func setLogger(cmd *cobra.Command, args []string) {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	switch logLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	ctx := context.WithValue(cmd.Context(), "logger", sugar)
	cmd.SetContext(ctx)
}
