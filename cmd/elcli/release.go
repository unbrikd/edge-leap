package elcli

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unbrikd/edge-leap/internal/azure"
	"github.com/unbrikd/edge-leap/internal/releaser"
	"github.com/unbrikd/edge-leap/internal/utils"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Handles the release of an application",
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile != "" {
			if _, err := loadConfig(); err != nil {
				fmt.Printf("error loading configuration: %v\n", err)
				os.Exit(1)
			}
		}
		executeRelease()
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	// Deployment configuration
	releaseCmd.Flags().StringVar(&config.Deployment.Id, "id", viper.GetString("deployment.id"), "id to use for deployment (must be kebab-case)")
	viper.BindPFlag("deployment.id", releaseCmd.Flags().Lookup("id"))

	releaseCmd.Flags().Int16VarP(&config.Deployment.Priority, "priority", "p", 50, "module deployment priority")
	viper.BindPFlag("deployment.priority", releaseCmd.Flags().Lookup("priority"))

	releaseCmd.Flags().StringVarP(&config.Deployment.TargetCondition, "target-condition", "t", viper.GetString("deployment.target-condition"), "target condition for the deployment")
	viper.BindPFlag("deployment.target-condition", releaseCmd.Flags().Lookup("target-condition"))

	// Device configuration
	releaseCmd.Flags().StringVar(&config.Device.Name, "device-name", viper.GetString("device.name"), "device name to deploy the module to")
	viper.BindPFlag("device.name", releaseCmd.Flags().Lookup("device-name"))

	// Module configuration
	releaseCmd.Flags().StringVarP(&config.Module.Name, "module-name", "m", viper.GetString("module.name"), "desired module name to show in the iotedge list (must be camelCase)")
	viper.BindPFlag("module.name", releaseCmd.Flags().Lookup("module-name"))

	releaseCmd.Flags().StringVar(&config.Module.CreateOptions, "create-options", viper.GetString("module.create-options"), "runtime settings for the container of the module (json string)")
	viper.BindPFlag("module.create-options", releaseCmd.Flags().Lookup("create-options"))

	releaseCmd.Flags().StringVarP(&config.Module.StartupOrder, "startup-order", "s", viper.GetString("module.startup-order"), "module startup order")
	viper.BindPFlag("module.startup-order", releaseCmd.Flags().Lookup("startup-order"))

	releaseCmd.Flags().StringVarP(&config.Module.Image, "image", "i", viper.GetString("module.image"), "module image URL (must be a valid docker image URL)")
	viper.BindPFlag("module.image", releaseCmd.Flags().Lookup("image"))

	releaseCmd.Flags().StringArrayVarP(&envFlag, "env", "e", nil, "environment variables for the module (key=value)")
	viper.BindPFlag("module.env", releaseCmd.Flags().Lookup("env"))

	// Infra configuration
	releaseCmd.Flags().StringVar(&config.Infra.Hub, "hub", "", "the name of the iot hub to send the deployment to")
	viper.BindPFlag("infra.hub", releaseCmd.Flags().Lookup("hub"))

	// Auth configuration
	releaseCmd.Flags().StringVar(&config.Auth.Token, "token", "", "token to authenticate the client")
	viper.BindPFlag("auth.token", releaseCmd.Flags().Lookup("token"))
}

// executeRelease handles the release of a module taking the configuration file or the flags.
// The flags have precedence over the configuration file.
func executeRelease() {
	moduleEnv, err := utils.StringArraySplitToMap(envFlag, "=")
	if err != nil {
		fmt.Printf("failed to parse environment variables: %v", err)
		os.Exit(1)
	}

	c := azure.NewClient(nil).WithAuthToken(config.Auth.Token)
	c.BaseURL, _ = url.Parse(fmt.Sprintf("https://%s.azure-devices.net/", config.Infra.Hub))

	d := azure.Configuration{
		Id:              config.Deployment.Id,
		Priority:        config.Deployment.Priority,
		TargetCondition: config.Deployment.TargetCondition,
	}
	d.SetContent(config.Module.Name, config.Module.Image, config.Module.CreateOptions, config.Module.StartupOrder, moduleEnv)

	fmt.Print(d)

	r := releaser.AzureReleaser{Client: c}
	err = r.ReleaseModule(&d)
	if err != nil {
		fmt.Printf("failed to release module: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", config.Deployment.Id)
}
