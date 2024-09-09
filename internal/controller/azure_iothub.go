package controller

import (
	"context"
)

type ConfigurationsService service

// GetLayeredDeployment gets a layered deployment from the Azure IoT Hub given its id. The function
// returns a LayeredManifest struct with the parsed deployment information, a nil pointer if the
// deployment was not found in the hub or an error if the api response is not 200.
// response is not 200.
func (s *ConfigurationsService) GetConfiguration(ctx context.Context, id string) (*LayeredManifest, error) {
	// u := fmt.Sprintf("configurations/%s?api-version=2021-04-12", id)

	req, err := s.client.NewRequest("GET", "https://iot-bems-dev-euw-p01.azure-devices.net/configurations/test-edge-leap?api-version=2021-04-12", nil)
	if err != nil {
		return nil, err
	}

	c := new(LayeredManifest)
	_, err = s.client.Do(req, c)
	if err != nil {
		return nil, err
	}

	return c, err

	// return c, nil

	// req, err := prepareRequest(http.MethodGet, url, az.token)
	// if err != nil {
	// 	return nil, fmt.Errorf("error creating request: %v", err)
	// }

	// res, err := apiCall(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("error calling api: %v", err)
	// }
	// defer res.Body.Close()

	// if res.StatusCode == http.StatusUnauthorized {
	// 	return nil, fmt.Errorf("unauthorized, please set a valid token")
	// }

	// if res.StatusCode == http.StatusNotFound {
	// 	return nil, nil
	// }

	// if res.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("unexpected hub behaviour: %v", res.Status)
	// }

	// lm := LayeredManifest{}
	// parser := json.NewDecoder(res.Body)
	// if err = parser.Decode(&lm); err != nil {
	// 	return nil, fmt.Errorf("failed to decode json response: %v", err)
	// }

}
