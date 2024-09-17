package drafter

import (
	"fmt"
	"net/http"

	"github.com/unbrikd/edge-leap/internal/azure"
)

type AzureDrafter struct {
	Client *azure.Client
}

func (d *AzureDrafter) SetOnDevice(deviceId, moduleName, moduleVersion string) (*azure.Twin, error) {
	twinTags := map[string]interface{}{
		"tags": map[string]interface{}{
			"application": map[string]string{
				moduleName: moduleVersion,
			},
		},
	}

	t, res, err := d.Client.Devices.UpdateTwinTags(deviceId, twinTags)
	if err != nil {
		return nil, err
	}

	if err = res.Expect(http.StatusOK); err != nil {
		// return nil, fmt.Errorf("failed to update the device twin: %v", res.Response.Header["Iothub-Errorcode"])
		return nil, fmt.Errorf("failed to update the device twin: %v", res.Response.Status)
	}

	return t, nil
}
