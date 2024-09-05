/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unbrikd/edge-leap/internal/configuration"
	"github.com/unbrikd/edge-leap/internal/controller"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Handles the deployment operations",
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.Help()
	},
}

var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "deploy the manifest as a layered deployment",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := configuration.Load(file)
		if err != nil {
			fmt.Printf("failed to load configuration: %v\n", err)
			return
		}

		ctl := controller.New(c, "azure")

		if err = ctl.DeleteLayeredDeployment(c); err != nil {
			fmt.Printf("failed to delete layered deployment: %v\n", err)
			return
		}

		if err = ctl.CreateLayeredDeployment(c); err != nil {
			fmt.Printf("failed to create layered deployment: %v\n", err)
			return
		}

		fmt.Printf("%s-%s\n", c.Image.Repo, c.Id)
	},
}

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Deploy the module in the device",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := configuration.Load(file)
		if err != nil {
			fmt.Printf("failed to load configuration: %v\n", err)
			return
		}

		ctl := controller.New(c, "azure")

		if err = ctl.UpdateDeviceTwin(c.Device.Name, map[string]string{c.Module.Name: c.Id}); err != nil {
			fmt.Printf("failed to update device twin: %v\n", err)
			return
		}

		fmt.Printf("%s: %s\n", c.Module.Name, c.Id)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.AddCommand(manifestCmd, moduleCmd)
}
