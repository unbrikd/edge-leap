package controller

import (
	"github.com/unbrikd/edge-leap/internal/configuration"
)

type Controller interface {
	// UpdateTwin updates the device twin for the given device.
	UpdateDeviceTwin(dn string, kv map[string]string) error

	// CreateLayeredDeployment creates a layered deployment for the given module.
	CreateLayeredDeployment(d Deployment) error

	// DeleteLayeredDeployment deletes a layered deployment for the given module.
	DeleteLayeredDeployment(id string) error

	// GetLayeredDeployment gets a layered deployment for the given module.
	GetLayeredDeployment(id string) (*LayeredManifest, error)
}

func New(c *configuration.Configuration, p string) Controller {
	switch p {
	case "Azure":
		return Azure(c.Infra.Hub, c.Auth.Token)
	default:
		return nil
	}
}
