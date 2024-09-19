package releaser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/unbrikd/edge-leap/internal/azure"
)

type AzureReleaser struct {
	// Azure client
	Client *azure.Client
}

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
