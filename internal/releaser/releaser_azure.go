package releaser

import (
	"context"

	"github.com/unbrikd/edge-leap/internal/azure"
)

type AzureReleaser struct {
	// Azure client
	Client *azure.Client
}

func (az *AzureReleaser) ReleaseModule(v *azure.Configuration) error {
	_, err := az.Client.Configurations.CreateConfiguration(context.Background(), *v)
	return err
}
