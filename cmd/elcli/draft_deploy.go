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
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := loadConfig(); err != nil {
			fmt.Printf("error loading configuration: %v\n", err)
			os.Exit(1)
		}
		preExecuteChecksDraftDeploy()
		executeDraftDeploy()
	},
}

func init() {
	draftCmd.AddCommand(draftDeployCmd)

	// Deployment configuration
	draftDeployCmd.Flags().StringVar(&config.Deployment.Id, "id", viper.GetString("deployment.id"), "id to use for deployment (must be kebab-case)")
	viper.BindPFlag("deployment.id", draftDeployCmd.Flags().Lookup("id"))

	draftDeployCmd.Flags().Int16VarP(&config.Deployment.Priority, "priority", "p", 50, "module deployment priority")
	viper.BindPFlag("deployment.priority", draftDeployCmd.Flags().Lookup("priority"))

	draftDeployCmd.Flags().StringVarP(&config.Deployment.TargetCondition, "target-condition", "t", viper.GetString("deployment.target-condition"), "target condition for the deployment")
	viper.BindPFlag("deployment.target-condition", draftDeployCmd.Flags().Lookup("target-condition"))

	// Device configuration
	draftDeployCmd.Flags().StringVar(&config.Device.Name, "device-name", viper.GetString("device.name"), "device name to deploy the module to")
	viper.BindPFlag("device.name", draftDeployCmd.Flags().Lookup("device-name"))

	// Module configuration
	draftDeployCmd.Flags().StringVarP(&config.Module.Name, "module-name", "m", viper.GetString("module.name"), "desired module name to show in the iotedge list (must be camelCase)")
	viper.BindPFlag("module.name", draftDeployCmd.Flags().Lookup("module-name"))

	draftDeployCmd.Flags().StringVar(&config.Module.CreateOptions, "create-options", viper.GetString("module.create-options"), "runtime settings for the container of the module (json string)")
	viper.BindPFlag("module.create-options", draftDeployCmd.Flags().Lookup("create-options"))

	draftDeployCmd.Flags().IntVarP(&config.Module.StartupOrder, "startup-order", "s", viper.GetInt("module.startup-order"), "module startup order")
	viper.BindPFlag("module.startup-order", draftDeployCmd.Flags().Lookup("startup-order"))

	draftDeployCmd.Flags().StringVarP(&config.Module.Image, "image", "i", viper.GetString("module.image"), "module image URL (must be a valid docker image URL)")
	viper.BindPFlag("module.image", draftDeployCmd.Flags().Lookup("image"))

	draftDeployCmd.Flags().StringSliceVarP(&config.Module.Env, "env", "e", nil, "environment variables for the module (key=value)")
	viper.BindPFlag("module.env", draftDeployCmd.Flags().Lookup("env"))

	// Infra configuration
	draftDeployCmd.Flags().StringVar(&config.Infra.Hub, "hub", "", "the name of the iot hub to send the deployment to")
	viper.BindPFlag("infra.hub", draftDeployCmd.Flags().Lookup("hub"))

	// Auth configuration
	draftDeployCmd.Flags().StringVar(&config.Auth.Token, "token", "", "token to authenticate the client")
	viper.BindPFlag("auth.token", draftDeployCmd.Flags().Lookup("token"))
}

// preExecuteChecksDraftDeploy checks if the required flags are set before executing the draft deploy command
func preExecuteChecksDraftDeploy() {
	for _, flag := range []string{"deployment.id", "device.name", "module.name"} {
		if viper.GetString(flag) == "" {
			fmt.Printf("error: %s is required\n", flag)
			os.Exit(1)
		}
	}
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
