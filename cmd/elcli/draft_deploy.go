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

var draftDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a draft module",
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
		checkRequired("deployment.id", "device.name", "module.name")
	},
	Run: func(cmd *cobra.Command, args []string) {
		executeDraftDeploy()
	},
}

func init() {
	// Deployment configuration
	draftDeployCmd.Flags().StringVar(&config.Deployment.Id, "id", "", "id to use for deployment (must be kebab-case)")
	draftDeployCmd.Flags().Int16VarP(&config.Deployment.Priority, "priority", "p", 50, "module deployment priority")
	draftDeployCmd.Flags().StringVarP(&config.Deployment.TargetCondition, "target-condition", "t", "", "target condition for the deployment")
	// Device configuration
	draftDeployCmd.Flags().StringVar(&config.Device.Name, "device-name", "", "device name to deploy the module to")
	// Module configuration
	draftDeployCmd.Flags().StringVarP(&config.Module.Name, "module-name", "m", "", "desired module name to show in the iotedge list (must be camelCase)")
	draftDeployCmd.Flags().StringVar(&config.Module.CreateOptions, "create-options", "", "runtime settings for the container of the module (json string)")
	draftDeployCmd.Flags().IntVarP(&config.Module.StartupOrder, "startup-order", "s", 20, "module startup order")
	draftDeployCmd.Flags().StringVarP(&config.Module.Image, "image", "i", "", "module image URL (must be a valid docker image URL)")
	draftDeployCmd.Flags().StringSliceVarP(&config.Module.Env, "env", "e", nil, "environment variables for the module (key=value)")
	// Infra configuration
	draftDeployCmd.Flags().StringVar(&config.Infra.Hub, "hub", "", "the name of the iot hub to send the deployment to")
	// Auth configuration
	draftDeployCmd.Flags().StringVar(&config.Auth.Token, "token", "", "token to authenticate the client")
}

func executeDraftDeploy() {
	c := azure.NewClient(nil).WithAuthToken(config.Auth.Token)
	c.BaseURL, _ = url.Parse(fmt.Sprintf("https://%s.azure-devices.net/", config.Infra.Hub))

	moduleEnv, err := utils.StringArraySplitToMap(config.Module.Env, "=")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := releaser.AzureReleaser{Client: c}
	if err := r.SetModuleOnDevice(config.Device.Name, config.Module.Name, config.Id); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d := azure.Configuration{
		Id:              fmt.Sprintf("%s-%s", config.Deployment.Id, config.Id),
		Priority:        config.Deployment.Priority,
		TargetCondition: fmt.Sprintf("tags.application.%s='%s'", config.Module.Name, config.Id),
	}
	d.SetContent(config.Module.Name, config.Module.Image, config.Module.CreateOptions, config.Module.StartupOrder, moduleEnv)

	if err := r.ReleaseModule(&d); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%s-%s\n", config.Deployment.Id, config.Id)
}
