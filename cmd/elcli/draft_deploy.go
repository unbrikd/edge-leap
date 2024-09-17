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

var draftDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a draft module",
	Run: func(cmd *cobra.Command, args []string) {
		LoadConfig()
		preExecuteChecksDraftDeploy()
		executeDraftDeploy()
	},
}

func init() {
	draftCmd.AddCommand(draftDeployCmd)

	draftDeployCmd.Flags().StringVar(&config.Deployment.Id, "id", "", "Deployment ID")
	viper.BindPFlag("deployment.id", draftDeployCmd.Flags().Lookup("id"))

	draftDeployCmd.Flags().StringVar(&config.Device.Name, "device-name", viper.GetString("device.name"), "Device name in the IoT Hub")
	viper.BindPFlag("device.name", draftDeployCmd.Flags().Lookup("device-name"))

	draftDeployCmd.Flags().StringVar(&config.Device.Name, "module-name", viper.GetString("module.name"), "Module to be drafted")
	viper.BindPFlag("module.name", draftDeployCmd.Flags().Lookup("module-name"))
}

// preExecuteChecksDraftDeploy checks if the required flags are set before executing the draft deploy command
func preExecuteChecksDraftDeploy() {
	if config.Deployment.Id == "" {
		fmt.Println("error: deployment ID is required, use --id flag or set it in the configuration file")
		os.Exit(1)
	}

	if config.Device.Name == "" {
		fmt.Println("error: device name is required, use --device-name flag or set it in the configuration file")
		os.Exit(1)
	}

	if config.Module.Name == "" {
		fmt.Println("error: module name is required, use --module-name flag or set it in the configuration file")
		os.Exit(1)
	}
}

func executeDraftDeploy() {
	c := azure.NewClient(nil).WithAuthToken(config.Auth.Token)
	c.BaseURL, _ = url.Parse(fmt.Sprintf("https://%s.azure-devices.net/", config.Infra.Hub))

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
	d.SetContent(config.Module.Name, config.Image.Repo, config.Module.CreateOptions, config.Module.StartupOrder)

	if err := r.ReleaseModule(&d); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
