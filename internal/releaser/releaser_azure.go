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
	fmt.Printf("Releasing module %s\n", c.Id)
	_, res, err := az.Client.Configurations.GetConfiguration(context.Background(), c.Id)
	if err != nil {
		return err
	}

	if err = res.Expect(200, 404); err != nil {
		return err
	}

	if res.Is(200) {
		fmt.Printf("Module %s already exists, deleting\n", c.Id)
		res, err = az.Client.Configurations.DeleteConfiguration(c.Id)
		if err != nil {
			return err
		}
		if err = res.Expect(http.StatusNoContent); err != nil {
			return err
		}
	}

	fmt.Printf("Creating module %s\n", c.Id)
	_, res, err = az.Client.Configurations.CreateConfiguration(context.Background(), *c)
	if err != nil {
		return err
	}

	if err = res.Expect(200); err != nil {
		return err
	}

	return nil
}
