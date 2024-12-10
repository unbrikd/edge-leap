package elcli

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
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
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("deployment.id", cmd.Flags().Lookup("id"))
		viper.BindPFlag("deployment.priority", cmd.Flags().Lookup("priority"))
		viper.BindPFlag("deployment.target-condition", cmd.Flags().Lookup("target-condition"))
		viper.BindPFlag("device.name", cmd.Flags().Lookup("device-name"))
		viper.BindPFlag("module.name", cmd.Flags().Lookup("module-name"))
		viper.BindPFlag("module.create-options", cmd.Flags().Lookup("create-options"))
		viper.BindPFlag("module.image", cmd.Flags().Lookup("image"))
		viper.BindPFlag("module.startup-order", cmd.Flags().Lookup("startup-order"))
		viper.BindPFlag("module.env", cmd.Flags().Lookup("env"))
		viper.BindPFlag("infra.hub", cmd.Flags().Lookup("hub"))
		viper.BindPFlag("auth.token", cmd.Flags().Lookup("token"))

		loadConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		executeRelease()
	},
}

func init() {
	// Deployment configuration
	releaseCmd.Flags().StringVar(&config.Deployment.Id, "id", "", "id to use for deployment (must be kebab-case)")
	releaseCmd.Flags().Int16VarP(&config.Deployment.Priority, "priority", "p", 50, "module deployment priority")
	releaseCmd.Flags().StringVarP(&config.Deployment.TargetCondition, "target-condition", "t", "", "target condition for the deployment")
	// Device configuration
	releaseCmd.Flags().StringVar(&config.Device.Name, "device-name", "", "device name to deploy the module to")
	// Module configuration
	releaseCmd.Flags().StringVarP(&config.Module.Name, "module-name", "m", "", "desired module name to show in the iotedge list (must be camelCase)")
	releaseCmd.Flags().StringVar(&config.Module.CreateOptions, "create-options", "", "runtime settings for the container of the module (json string)")
	releaseCmd.Flags().IntVarP(&config.Module.StartupOrder, "startup-order", "s", 20, "module startup order")
	releaseCmd.Flags().StringVarP(&config.Module.Image, "image", "i", "", "module image URL (must be a valid docker image URL)")
	releaseCmd.Flags().StringSliceVarP(&config.Module.Env, "env", "e", nil, "environment variables for the module (key=value)")
	// Infra configuration
	releaseCmd.Flags().StringVar(&config.Infra.Hub, "hub", "", "the name of the iot hub to send the deployment to")
	// Auth configuration
	releaseCmd.Flags().StringVar(&config.Auth.Token, "token", "", "token to authenticate the client")
}

// executeRelease handles the release of a module taking the configuration file or the flags.
// The flags have precedence over the configuration file.
func executeRelease() {
	moduleEnv, err := utils.StringArraySplitToMap(config.Module.Env, "=")
	if err != nil {
		fmt.Printf("failed to parse environment variables: %v", err)
		os.Exit(1)
	}

	c := azure.NewClient(nil).WithAuthToken(config.Auth.Token)
	c.BaseURL, _ = url.Parse(fmt.Sprintf("https://%s.azure-devices.net/", config.Infra.Hub))

	releaseId := strings.Split(uuid.New().String(), "-")[4]
	d := azure.Configuration{
		Id:              config.Deployment.Id,
		Priority:        config.Deployment.Priority,
		TargetCondition: config.Deployment.TargetCondition,
		Labels: map[string]string{
			"releaseId": releaseId},
	}
	d.SetContent(config.Module.Name, config.Module.Image, config.Module.CreateOptions, config.Module.StartupOrder, moduleEnv)

	r := releaser.Azure(c)
	err = r.ReleaseModule(&d)
	if err != nil {
		fmt.Printf("failed to release module: %v", err)
		os.Exit(1)
	}

	fmt.Printf("%s, release %s\n", config.Deployment.Id, releaseId)
}
