package configuration

import (
	"fmt"
	"os"

	"github.com/unbrikd/edge-leap/internal/utils"
	"gopkg.in/yaml.v3"
)

const CONFIG_VERSION = 1

type Configuration struct {
	// Id is the unique identifier of the session.
	Id string `yaml:"id"`
	// Config is the path to the configuration file.
	Config string `yaml:"config"`
	// Version is the version of the configuration file.
	Version int `yaml:"version"`

	// Module struct holds the module information.
	Module struct {
		// Name is the name of the module in the edge workload controller runtime.
		Name string `yaml:"name"`
		// Version is the version of the module.
		Version string `yaml:"version,omitempty"`
		// Manifest is the path to the module manifest file.
		Manifest string `yaml:"manifest"`
	} `yaml:"module"`

	// Image struct holds the image information.
	Image struct {
		// Repo is the repository of the image to be pushed to the registry.
		Repo string `yaml:"repo"`
	} `yaml:"image"`

	// Device struct holds the development device information.
	Device struct {
		// Name is the name of the device in the cloud provider.
		Name string `yaml:"name"`
		// Arch is the architecture of the device.
		Arch string `yaml:"arch"`
	} `yaml:"device"`

	// Infra struct holds the infrastructure information.
	Infra struct {
		// Hub is the name of the IoT Hub where the development device is connected.
		Hub string `yaml:"hub,omitempty"`
		// Registry is the name of the container registry to push the images.
		Registry string `yaml:"registry,omitempty"`
	} `yaml:"infra"`

	Auth struct {
		Token string `yaml:"token,omitempty"`
	} `yaml:"auth"`
}

func New(cp string) (*Configuration, error) {
	c := &Configuration{
		Id:      utils.ShortUuid(),
		Config:  cp,
		Version: CONFIG_VERSION,
	}

	if err := c.toFile(c.Config); err != nil {
		return nil, err
	}

	return c, nil
}

// toFile creates a new configuration file with the session id, configuration path and
// version as content. The rest of the fields are empty to be filled by the user.
func (c *Configuration) toFile(path string) error {
	// Create a new file with the session id as content.
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create configuration file: %v", err)
	}
	defer f.Close()

	// Write the configuration struct to a yaml file.
	out, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("could not marshal configuration to yaml: %v", err)
	}

	_, err = f.Write(out)
	if err != nil {
		return fmt.Errorf("could not write configuration to configuration file: %v", err)
	}

	return nil
}

func Load(path string) (*Configuration, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file does not exist: %v", err)
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read configuration file: %v", err)
	}

	c := &Configuration{}
	if err := yaml.Unmarshal(f, c); err != nil {
		return nil, fmt.Errorf("could not unmarshal configuration file: %v", err)
	}

	return c, nil
}
