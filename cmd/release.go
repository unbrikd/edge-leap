package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/unbrikd/edge-leap/internal/azure"
	"github.com/unbrikd/edge-leap/internal/configuration"
)

var manifest, image, createOpts, startupOrder, targetCondition, token, module, id string
var priority int16

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Handles the release of an application",
	Run: func(cmd *cobra.Command, args []string) {
		executeRelease()
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	// Path to the manifest file
	releaseCmd.Flags().StringVarP(&manifest, "manifest", "m", "", "Path to the manifest file")
	// Docker image to set in the manifest
	releaseCmd.Flags().StringVarP(&image, "image", "i", "", "Docker image to set in the manifest")
	// Create options to set how the module is initialized from iotedge
	releaseCmd.Flags().StringVarP(&createOpts, "create-options", "c", "", "Options to set how the module is initialized from iotedge")
	// Priority of the module
	releaseCmd.Flags().Int16VarP(&priority, "priority", "p", 10, "Priority of the module")
	// Target condition to set in the manifest
	releaseCmd.Flags().StringVarP(&targetCondition, "target-condition", "t", "", "Target condition to set in the manifest")
	// Statup order of the module
	releaseCmd.Flags().StringVarP(&startupOrder, "startup-order", "s", "50", "Startup order of the module")
	// Token to authenticate with the azure hub
	releaseCmd.Flags().StringVar(&token, "token", "", "Token to authenticate with the azure hub")
	// Module name to release
	releaseCmd.Flags().StringVarP(&module, "module", "m", "", "Module name as shown in iotedge list command")
	// Deployment id to release
	releaseCmd.Flags().StringVar(&id, "id", "", "Deployment id to release")
}

func executeRelease() {
	cfg, err := configuration.Load(file)
	if err != nil {
		fmt.Printf("failed to load configuration: %v", err)
		return
	}

	if token != "" {
		cfg.Auth.Token = token
	}

	c := azure.NewClient(nil)
	c.WithAuthToken(token)
	c.BaseURL, _ = url.Parse(fmt.Sprintf("https://%s.azure-devices.net/", cfg.Infra.Hub))

	// m, err := c.Configurations.GetConfiguration(context.Background(), deploymentId)
	// if err != nil {
	// 	fmt.Printf("failed to get layered deployment: %v", err)
	// 	return
	// }

	// fmt.Println(m)

	// err = c.Configurations.DeleteConfiguration(deploymentId)
	// if err != nil {
	// 	fmt.Printf("failed to delete layered deployment: %v", err)
	// 	return
	// }

	d := azure.Configuration{
		Id:              id,
		Priority:        priority,
		TargetCondition: targetCondition,
	}
	d.SetContent(module, image, createOpts, startupOrder)

	fmt.Print(d)

	n, err := c.Configurations.CreateConfiguration(context.Background(), d)
	if err != nil {
		fmt.Printf("failed to create configuration: %v", err)
		return
	}

	fmt.Printf("created configuration: %v", n)
}
