package configuration

const CONFIG_VERSION = 1

type Configuration struct {
	// Id is the unique identifier of the session.
	Id string `mapstructure:"session"`
	// Version is the version of the configuration file.
	Version int `mapstructure:"version"`

	// Module struct holds the module information.
	Module struct {
		// Name is the name of the module in the edge workload controller runtime.
		Name string `mapstructure:"name,omitempty"`
		// StartupOrder is the startup order of the module in the cloud provider.
		StartupOrder string `mapstructure:"startup-order,omitempty"`
		// CreateOptions is the create options of the module in the cloud provider.
		CreateOptions string `mapstructure:"create-options,omitempty"`
	} `mapstructure:"module"`

	Deployment struct {
		// Id is the deployment id of the module in the cloud provider.
		Id string `mapstructure:"id"`
		// Priority is the priority of the module in the cloud provider.
		Priority int16 `mapstructure:"priority"`
		// TargetCondition is the target condition of the module in the cloud provider.
		TargetCondition string `mapstructure:"target-condition"`
	} `mapstructure:"deployment"`

	// Image struct holds the image information.
	Image struct {
		// Repo is the repository of the image to be pushed to the registry.
		Repo string `mapstructure:"repo"`
		// Tag is the tag of the image to be pushed to/from the registry.
		Tag string `mapstructure:"tag"`
	} `mapstructure:"image"`

	// Device struct holds the development device information.
	Device struct {
		// Name is the name of the device in the cloud provider.
		Name string `mapstructure:"name"`
		// Arch is the architecture of the device.
		Arch string `mapstructure:"arch"`
	} `mapstructure:"device"`

	// Infra struct holds the infrastructure information.
	Infra struct {
		// Hub is the name of the IoT Hub where the development device is connected.
		Hub string `mapstructure:"hub"`
		// Registry is the name of the container registry to push the images.
		Registry string `mapstructure:"registry"`
	} `mapstructure:"infra"`

	Auth struct {
		Token string `mapstructure:"token"`
	} `mapstructure:"auth"`
}
