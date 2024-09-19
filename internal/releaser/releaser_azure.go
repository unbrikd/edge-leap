package releaser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/unbrikd/edge-leap/internal/azure"
)

// AzureReleaser is the handler for releasing configurations to Azure IoT Hub and deploy modules to devices
type AzureReleaser struct {
	// Azure client
	Client *azure.Client
}

// ReleaseModule releases a new configuration to Azure IoT Hub. Devices configurations are left untouched.
// If a configuration with the same id already exists, it will be deleted and replaced.
// If the configuration fails to be created, the previous configuration will be restored.
func (az *AzureReleaser) ReleaseModule(c *azure.Configuration) error {
	currentConfig, err := az.configurationExists(c.Id)
	if err != nil {
		return err
	}

	if currentConfig != nil {
		err = az.configurationAttemptDelete(c.Id)
		if err != nil {
			return err
		}
	}

	err = az.configurationAttemptCreate(c)
	if err != nil {
		if currentConfig != nil {
			az.configurationAttemptCreate(currentConfig)
		}

		return err
	}

	return nil
}

// SetModuleOnDevice sets the module name and version on the device twin tags in order to allow drafts to be deployed without manual intervention.
// By convention the edge leap expects all layered deployments to have a target condition such as: "tags.application.<module_name> = '<module_version>'"
// If the convention is not followed, the deployment will be available but not applied to the device.
func (az *AzureReleaser) SetModuleOnDevice(deviceId, moduleName, moduleVersion string) error {
	twinTags := map[string]interface{}{
		"tags": map[string]interface{}{
			"application": map[string]string{
				moduleName: moduleVersion,
			},
		},
	}

	_, res, err := az.Client.Devices.UpdateTwinTags(deviceId, twinTags)
	if err != nil {
		return err
	}

	if err = res.Expect(http.StatusOK); err != nil {
		return fmt.Errorf("failed to update the device twin: %v", res.Response.Header["Iothub-Errorcode"])
	}

	return nil
}

// configurationExists checks if a configuration with the given id exists and returns it as a Configuration object.
// If the configuration does not exist, nil is returned.
func (az *AzureReleaser) configurationExists(id string) (*azure.Configuration, error) {
	c, res, err := az.Client.Configurations.GetConfiguration(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if err = res.Expect(http.StatusOK, http.StatusNotFound); err != nil {
		return nil, fmt.Errorf("failed to check configuration: %v", res.Response.Header["Iothub-Errorcode"])
	}

	if res.Is(http.StatusNotFound) {
		return nil, nil
	}

	return c, nil
}

// configurationAttemptCreate attempts to create a new configuration in Azure IoT Hub.
// In case of any error, the configuration will not be created and an error will be returned.
func (az *AzureReleaser) configurationAttemptCreate(c *azure.Configuration) error {
	_, res, err := az.Client.Configurations.CreateConfiguration(context.Background(), *c)
	if err != nil {
		return err
	}

	if err = res.Expect(http.StatusOK); err != nil {
		return fmt.Errorf("failed to create configuration: %v", res.Response.Header["Iothub-Errorcode"])
	}

	return nil
}

// configurationAttemptDelete attempts to delete a configuration in Azure IoT Hub.
// In case of any error, the configuration will not be deleted and an error will be returned.
func (az *AzureReleaser) configurationAttemptDelete(id string) error {
	res, err := az.Client.Configurations.DeleteConfiguration(id)
	if err != nil {
		return err
	}

	if err = res.Expect(http.StatusNoContent); err != nil {
		return err
	}

	return nil
}
