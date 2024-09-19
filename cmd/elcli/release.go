package elcli

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unbrikd/edge-leap/internal/azure"
	"github.com/unbrikd/edge-leap/internal/releaser"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Handles the release of an application",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := loadConfig(); err != nil {
			fmt.Printf("error loading configuration: %v\n", err)
			os.Exit(1)
		}
		executeRelease()
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	// Deployment configuration
	releaseCmd.Flags().StringVar(&config.Deployment.Id, "id", viper.GetString("deployment.id"), "Deployment id to release")
	viper.BindPFlag("deployment.id", releaseCmd.Flags().Lookup("id"))

	releaseCmd.Flags().Int16VarP(&config.Deployment.Priority, "priority", "p", 50, "Priority of the module")
	viper.BindPFlag("deployment.priority", releaseCmd.Flags().Lookup("priority"))

	releaseCmd.Flags().StringVarP(&config.Deployment.TargetCondition, "target-condition", "t", viper.GetString("deployment.target-condition"), "Target condition to set in the manifest")
	viper.BindPFlag("deployment.target-condition", releaseCmd.Flags().Lookup("target-condition"))

	// Device configuration
	releaseCmd.Flags().StringVar(&config.Device.Name, "device-name", viper.GetString("device.name"), "Device name in the IoT Hub")
	viper.BindPFlag("device.name", releaseCmd.Flags().Lookup("device-name"))

	// Module configuration
	releaseCmd.Flags().StringVarP(&config.Module.Name, "module-name", "m", viper.GetString("module.name"), "Module name as shown in iotedge list command")
	viper.BindPFlag("module.name", releaseCmd.Flags().Lookup("module-name"))

	releaseCmd.Flags().StringVar(&config.Module.CreateOptions, "create-options", viper.GetString("module.create-options"), "Options to set how the module is initialized from iotedge")
	viper.BindPFlag("module.create-options", releaseCmd.Flags().Lookup("create-options"))

	releaseCmd.Flags().StringVarP(&config.Module.StartupOrder, "startup-order", "s", viper.GetString("module.startup-order"), "Startup order of the module")
	viper.BindPFlag("module.startup-order", releaseCmd.Flags().Lookup("startup-order"))

	releaseCmd.Flags().StringVarP(&config.Module.Image, "image", "i", viper.GetString("module.image"), "Startup order of the module")
	viper.BindPFlag("module.image", releaseCmd.Flags().Lookup("image"))

	// Infra configuration
	releaseCmd.Flags().StringVar(&config.Infra.Hub, "hub", "", "IoT Hub name")
	viper.BindPFlag("infra.hub", releaseCmd.Flags().Lookup("hub"))

	// Auth configuration
	releaseCmd.Flags().StringVar(&config.Auth.Token, "token", "", "Token to authenticate the client")
	viper.BindPFlag("auth.token", releaseCmd.Flags().Lookup("token"))
}

// executeRelease handles the release of a module taking the configuration file or the flags.
// The flags have precedence over the configuration file.
func executeRelease() {
	c := azure.NewClient(nil).WithAuthToken(config.Auth.Token)
	c.BaseURL, _ = url.Parse(fmt.Sprintf("https://%s.azure-devices.net/", config.Infra.Hub))

	d := azure.Configuration{
		Id:              config.Deployment.Id,
		Priority:        config.Deployment.Priority,
		TargetCondition: config.Deployment.TargetCondition,
	}
	d.SetContent(config.Module.Name, config.Module.Image, config.Module.CreateOptions, config.Module.StartupOrder)

	r := releaser.AzureReleaser{Client: c}
	err := r.ReleaseModule(&d)
	if err != nil {
		fmt.Printf("failed to release module: %v", err)
		return
	}

	fmt.Printf("%s\n", config.Deployment.Id)
}
