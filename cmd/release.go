package cmd

import (
	"fmt"
	"net/url"

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
		initConfig()
		executeRelease()
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	// Module configuration
	releaseCmd.Flags().StringVarP(&config.Deployment.CreateOptions, "create-options", "c", "", "Options to set how the module is initialized from iotedge")
	viper.BindPFlag("module.createOptions", releaseCmd.Flags().Lookup("create-options"))

	releaseCmd.Flags().StringVarP(&config.Module.Name, "module-name", "m", "", "Module name as shown in iotedge list command")
	viper.BindPFlag("module.name", releaseCmd.Flags().Lookup("module.name"))

	// Deployment configuration
	releaseCmd.Flags().StringVar(&config.Deployment.Id, "id", "", "Deployment id to release")
	viper.BindPFlag("deployment.id", releaseCmd.Flags().Lookup("id"))

	releaseCmd.Flags().StringVar(&config.Deployment.Manifest, "json-file", "", "Path to the manifest file")
	viper.BindPFlag("deployment.manifest", releaseCmd.Flags().Lookup("json-file"))

	releaseCmd.Flags().Int16VarP(&config.Deployment.Priority, "priority", "p", 10, "Priority of the module")
	viper.BindPFlag("deployment.priority", releaseCmd.Flags().Lookup("priority"))

	releaseCmd.Flags().StringVarP(&config.Deployment.TargetCondition, "target-condition", "t", "", "Target condition to set in the manifest")
	viper.BindPFlag("deployment.targetCondition", releaseCmd.Flags().Lookup("target-condition"))

	releaseCmd.Flags().StringVarP(&config.Deployment.StartupOrder, "startup-order", "s", "50", "Startup order of the module")
	viper.BindPFlag("deployment.startupOrder", releaseCmd.Flags().Lookup("startup-order"))

	// Image configuration
	releaseCmd.Flags().StringVarP(&config.Image.Repo, "image", "i", "", "Docker image to set in the manifest")
	viper.BindPFlag("image.repo", releaseCmd.Flags().Lookup("image"))

	releaseCmd.Flags().StringVar(&config.Image.Tag, "tag", "latest", "Docker image tag to set in the manifest")
	viper.BindPFlag("image.tag", releaseCmd.Flags().Lookup("tag"))

	// Infra configuration
	releaseCmd.Flags().StringVar(&config.Infra.Hub, "hub", "", "IoT Hub name")
	viper.BindPFlag("infra.hub", releaseCmd.Flags().Lookup("hub"))

	// Auth configuration
	releaseCmd.Flags().StringVar(&config.Auth.Token, "token", "", "Token to authenticate the client")
	viper.BindPFlag("auth.token", releaseCmd.Flags().Lookup("token"))
}

func executeRelease() {
	c := azure.NewClient(nil).WithAuthToken(config.Auth.Token)
	c.BaseURL, _ = url.Parse(fmt.Sprintf("https://%s.azure-devices.net/", config.Infra.Hub))

	d := azure.Configuration{
		Id:              config.Deployment.Id,
		Priority:        config.Deployment.Priority,
		TargetCondition: config.Deployment.TargetCondition,
	}
	d.SetContent(config.Module.Name, config.Image.Repo, config.Deployment.CreateOptions, config.Deployment.StartupOrder)

	r := releaser.AzureReleaser{Client: c}
	err := r.ReleaseModule(&d)
	if err != nil {
		fmt.Printf("failed to release module: %v", err)
		return
	}
}
