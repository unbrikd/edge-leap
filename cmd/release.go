package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/unbrikd/edge-leap/internal/configuration"
	"github.com/unbrikd/edge-leap/internal/controller"
)

var manifest, image, createOpts, startupOrder, targetCondition, token, deploymentId string
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
	// releaseCmd.MarkFlagRequired("manifest")

	// Docker image to set in the manifest
	releaseCmd.Flags().StringVarP(&image, "image", "i", "", "Docker image to set in the manifest")
	// releaseCmd.MarkFlagRequired("image")

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
	releaseCmd.Flags().StringVar(&deploymentId, "id", "", "Id for the module deployment manifest")
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

	c := controller.NewClient(nil)
	c.WithAuthToken(token)
	c.BaseURL, _ = url.Parse(fmt.Sprintf("https://%s.azure-devices.net/", cfg.Infra.Hub))

	m, err := c.Configurations.GetConfiguration(context.Background(), deploymentId)
	if err != nil {
		fmt.Printf("failed to get layered deployment: %v", err)
		return
	}

	fmt.Println(m)

	// ctlr := controller.New(cfg, "Azure")
	// rel := releaser.New(ctlr)

	// d := controller.Deployment{
	// 	Id:              deploymentId,
	// 	Priority:        priority,
	// 	TargetCondition: targetCondition,
	// 	ManifestPath:    manifest,
	// }

	// if err = rel.ReleaseModule(d); err != nil {
	// 	fmt.Printf("failed to release module: %v", err)
	// 	return
	// }
}
