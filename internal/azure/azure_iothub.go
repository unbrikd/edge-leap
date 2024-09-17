package azure

import (
	"context"
	"fmt"
)

type ConfigurationsService service
type DevicesService service

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

type Twin struct {
	DeviceId string      `json:"deviceId"`
	Tags     interface{} `json:"tags"`
}

// GetConfiguration retrieves a configuration from the Azure IoT Hub. A configuration object is returned if the
// operation is successful, otherwise an error is returned and the configuration object is nil.
func (s *ConfigurationsService) GetConfiguration(ctx context.Context, id string) (*Configuration, *Response, error) {
	u := fmt.Sprintf("configurations/%s?api-version=2021-04-12", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	c := new(Configuration)
	res, err := s.client.Do(req, c)
	if err != nil {
		return nil, nil, err
	}

	return c, &Response{res}, err
}

func (s *ConfigurationsService) CreateConfiguration(ctx context.Context, c Configuration) (*Configuration, *Response, error) {
	u := fmt.Sprintf("configurations/%s?api-version=2021-04-12", c.Id)

	req, err := s.client.NewRequest("PUT", u, c)
	if err != nil {
		return nil, nil, err
	}

	cNew := new(Configuration)
	res, err := s.client.Do(req, cNew)
	if err != nil {
		return nil, nil, err
	}

	return cNew, &Response{res}, nil
}

// DeleteConfiguration deletes a configuration from the Azure IoT Hub. An error is returned if the operation is not
// successful.
func (s *ConfigurationsService) DeleteConfiguration(id string) (*Response, error) {
	u := fmt.Sprintf("configurations/%s?api-version=2021-04-12", id)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return &Response{res}, nil
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

func (d *DevicesService) GetTwin(deviceId string) (*Twin, *Response, error) {
	u := fmt.Sprintf("twins/%s?api-version=2021-04-12", deviceId)

	req, err := d.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	t := new(Twin)
	res, err := d.client.Do(req, t)
	if err != nil {
		return nil, nil, err
	}

	return t, &Response{res}, nil
}

func (d *DevicesService) UpdateTwinTags(deviceId string, tags map[string]interface{}) (*Twin, *Response, error) {
	u := fmt.Sprintf("twins/%s?api-version=2021-04-12", deviceId)

	req, err := d.client.NewRequest("PATCH", u, tags)
	if err != nil {
		return nil, nil, err
	}

	tNew := new(Twin)
	res, err := d.client.Do(req, tNew)
	if err != nil {
		return nil, nil, err
	}

	return tNew, &Response{res}, nil
}
