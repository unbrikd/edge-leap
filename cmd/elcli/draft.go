package elcli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var draftCmd = &cobra.Command{
	Use:   "draft",
	Short: "Start a new module draft session",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
	},
	Run: func(cmd *cobra.Command, args []string) {
		executeDraft()
	},
}

func init() {
	rootCmd.AddCommand(draftCmd)
	draftCmd.AddCommand(draftDeployCmd)
	draftCmd.AddCommand(draftNewCmd)

	// Deployment configuration
	draftCmd.PersistentFlags().StringVar(&config.Deployment.Id, "id", "", "id to use for deployment (must be kebab-case)")
	draftCmd.PersistentFlags().Int16VarP(&config.Deployment.Priority, "priority", "p", 50, "module deployment priority")
	draftCmd.PersistentFlags().StringVarP(&config.Deployment.TargetCondition, "target-condition", "t", "", "target condition for the deployment")
	// Device configuration
	draftCmd.PersistentFlags().StringVar(&config.Device.Name, "device-name", "", "device name to deploy the module to")
	// Module configuration
	draftCmd.PersistentFlags().StringVarP(&config.Module.Name, "module-name", "m", "", "desired module name to show in the iotedge list (must be camelCase)")
	draftCmd.PersistentFlags().StringVar(&config.Module.CreateOptions, "create-options", "", "runtime settings for the container of the module (json string)")
	draftCmd.PersistentFlags().IntVarP(&config.Module.StartupOrder, "startup-order", "s", 20, "module startup order")
	draftCmd.PersistentFlags().StringVarP(&config.Module.Image, "image", "i", "", "module image URL (must be a valid docker image URL)")
	draftCmd.PersistentFlags().StringSliceVarP(&config.Module.Env, "env", "e", nil, "environment variables for the module (key=value)")
	// Infra configuration
	draftCmd.PersistentFlags().StringVar(&config.Infra.Hub, "hub", "", "the name of the iot hub to send the deployment to")
	// Auth configuration
	draftCmd.PersistentFlags().StringVar(&config.Auth.Token, "token", "", "token to authenticate the client")
}

func executeDraft() {
	rootCmd.Help()
}
