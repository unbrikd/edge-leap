package azure

import (
	"context"
	"fmt"
	"net/http"
)

type ConfigurationsService service

// Configuration represents an Azure IoT Hub configuration. This schema is described at:
// https://learn.microsoft.com/en-us/rest/api/iothub/service/configuration/get?view=rest-iothub-service-2021-11-30
type Configuration struct {
	Content            map[string]interface{} `json:"content"`
	CreatedTimeUtc     string                 `json:"createdTimeUtc,omitempty"`
	ETag               string                 `json:"etag,omitempty"`
	Id                 string                 `json:"id"`
	Labels             map[string]string      `json:"labels,omitempty"`
	LastUpdatedTimeUtc string                 `json:"lastUpdatedTimeUtc,omitempty"`
	Metrics            interface{}            `json:"metrics,omitempty"`
	Priority           int16                  `json:"priority"`
	SchemaVersion      string                 `json:"schemaVersion,omitempty"`
	SystemMetrics      interface{}            `json:"systemMetrics,omitempty"`
	TargetCondition    string                 `json:"targetCondition"`
}

// GetConfiguration retrieves a configuration from the Azure IoT Hub. A configuration object is returned if the
// operation is successful, otherwise an error is returned and the configuration object is nil.
func (s *ConfigurationsService) GetConfiguration(ctx context.Context, id string) (*Configuration, error) {
	u := fmt.Sprintf("configurations/%s?api-version=2021-04-12", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	c := new(Configuration)
	_, err = s.client.Do(req, c)
	if err != nil {
		return nil, err
	}

	return c, err
}

func (s *ConfigurationsService) CreateConfiguration(ctx context.Context, c Configuration) (*Configuration, error) {
	u := fmt.Sprintf("configurations/%s?api-version=2021-04-12", c.Id)

	req, err := s.client.NewRequest("PUT", u, c)
	if err != nil {
		return nil, err
	}

	cNew := new(Configuration)
	_, err = s.client.Do(req, cNew)
	if err != nil {
		return nil, err
	}

	return cNew, err
}

// DeleteConfiguration deletes a configuration from the Azure IoT Hub. An error is returned if the operation is not
// successful.
func (s *ConfigurationsService) DeleteConfiguration(id string) error {
	u := fmt.Sprintf("configurations/%s?api-version=2021-04-12", id)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	res, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized, please set a valid token")
	}

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected hub behaviour: %v", res.Status)
	}

	return nil
}

func (c *Configuration) SetContent(mod, img, opts, so string) {
	props := fmt.Sprintf("properties.desired.%s", mod)
	contents := map[string]interface{}{
		"modulesContent": map[string]interface{}{
			"$edgeAgent": map[string]interface{}{
				props: map[string]interface{}{
					"settings": map[string]string{
						"image":         img,
						"createOptions": opts,
					},
					"startupOrder": so,
				},
			},
		},
	}

	c.Content = contents
}
