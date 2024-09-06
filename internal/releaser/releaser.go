package releaser

import (
	"fmt"

	"github.com/unbrikd/edge-leap/internal/controller"
)

type Releaser struct {
	controller controller.Controller
}

func New(c controller.Controller) *Releaser {
	return &Releaser{
		controller: c,
	}
}

func (r *Releaser) ReleaseModule(d controller.Deployment) error {
	// check if any other deployment is available
	search, err := r.controller.GetLayeredDeployment(d.Id)
	if err != nil {
		return fmt.Errorf("failed to get layered deployment: %w", err)
	}

	// if layered deployment exists then we need to delete it first
	if search != nil {
		if err = r.controller.DeleteLayeredDeployment(d.Id); err != nil {
			return fmt.Errorf("failed to delete layered deployment: %w", err)
		}
	}

	if err = r.controller.CreateLayeredDeployment(d); err != nil {
		return fmt.Errorf("failed to create layered deployment: %w", err)
	}

	return nil
}
