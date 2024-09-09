package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// LayeredManifest is the structure that represents a manifest for a layered deployment
type LayeredManifest struct {
	// Id is the unique identifier of the layered deployment
	Id string `json:"id"`
	// Priority is the priority of the layered deployment
	Priority int16 `json:"priority"`
	// TargetCondition sets the in which conditions the deployment should be applied
	TargetCondition string `json:"targetCondition"`
	// Content sets the module to be deployed and its settings
	Content struct {
		ModulesContent struct {
			EdgeAgent map[string]interface{} `json:"$edgeAgent"`
		} `json:"modulesContent"`
	} `json:"content"`
}

// Deployment is the structure that aggregates the information needed to create a deployment
type Deployment struct {
	ManifestPath    string
	Id              string
	Priority        int16
	TargetCondition string
}

// azCtl is a controller for Azure IoT Hub.
type azCtl struct {
	// hubName is the name of the Azure IoT Hub.
	hubName string
	// token is the authentication token for the Azure IoT Hub.
	token string
}

// Azure initializes and returns a new Azure IoT Hub controller.
func Azure(hn string, t string) *azCtl {
	return &azCtl{
		hubName: hn,
		token:   t,
	}
}

func (az *azCtl) UpdateDeviceTwin(dn string, kv map[string]string) error {
	url := fmt.Sprintf("https://%s.azure-devices.net/twins/%s?api-version=2021-04-12", az.hubName, dn)

	req, err := prepareRequest("PATCH", url, az.token)
	if err != nil {
		return err
	}

	for k, v := range kv {
		newTags := fmt.Sprintf(`{"tags": {"application": {"%s": "%s"}}`, k, v)
		req.Body = io.NopCloser(strings.NewReader(newTags))

		res, err := apiCall(req)
		if err != nil {
			return err
		}

		if res.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("unauthorized, please set a valid token")
		}

		if res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("device %s not found", dn)
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected hub behaviour: %v", res.Status)
		}
	}

	return nil
}

// CreateLayeredDeployment creates a new layered deployment in the Azure IoT Hub. The deployment
// manifest is formed from the given manifest file and it is included the given layered deployment
// options (id, priority, target condition). The function returns an error if the api response
// status is not 200.
func (az *azCtl) CreateLayeredDeployment(d Deployment) error {
	url := fmt.Sprintf("https://%s.azure-devices.net/configurations/%s?api-version=2021-04-12", az.hubName, d.Id)

	contents, _ := os.ReadFile(d.ManifestPath)

	lm := &LayeredManifest{}
	if err := json.Unmarshal(contents, lm); err != nil {
		return fmt.Errorf("cannot unmarshal manifest: %v", err)
	}
	lm.Id = d.Id
	lm.Priority = d.Priority
	lm.TargetCondition = d.TargetCondition

	manifest, _ := json.Marshal(lm)

	// create the request
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", az.token)
	req.Header.Set("Content-Type", "application/json")

	req.Body = io.NopCloser(bytes.NewReader(manifest))

	res, err := apiCall(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized, please set a valid token")
	}

	if res.StatusCode == http.StatusConflict {
		return fmt.Errorf("deployment %s already exists", d.Id)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected hub behaviour: %v", res.Status)
	}

	return nil
}

// DeleteLayeredDeployment deletes a layered deployment from the Azure IoT Hub given its id. The
// function returns an error if the api response status is not 204 or nil if the deployment was
// successfully deleted.
func (az *azCtl) DeleteLayeredDeployment(id string) error {
	url := fmt.Sprintf("https://%s.azure-devices.net/configurations/%s?api-version=2021-04-12", az.hubName, id)

	req, err := prepareRequest(http.MethodDelete, url, az.token)
	if err != nil {
		return err
	}

	res, err := apiCall(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized, please set a valid token")
	}

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected hub behaviour: %v", res.Status)
	}

	return nil
}

// GetLayeredDeployment gets a layered deployment from the Azure IoT Hub given its id. The function
// returns a LayeredManifest struct with the parsed deployment information, a nil pointer if the
// deployment was not found in the hub or an error if the api response is not 200.
// response is not 200.
func (az *azCtl) GetLayeredDeployment(id string) (*LayeredManifest, error) {
	url := fmt.Sprintf("https://%s.azure-devices.net/configurations/%s?api-version=2021-04-12", az.hubName, id)

	req, err := prepareRequest(http.MethodGet, url, az.token)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	res, err := apiCall(req)
	if err != nil {
		return nil, fmt.Errorf("error calling api: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized, please set a valid token")
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected hub behaviour: %v", res.Status)
	}

	lm := LayeredManifest{}
	parser := json.NewDecoder(res.Body)
	if err = parser.Decode(&lm); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %v", err)
	}

	return &lm, nil
}

// A Client manages communication with the GitHub API.
type Client struct {
	client  *http.Client // HTTP client used to communicate with the API.
	BaseURL *url.URL

	Configurations *ConfigurationsService
}

type service struct {
	client *Client

	// Base URL for API requests. Every service should have a base URL set, since the azure API base URL differs from
	// resource to resource.
	baseURL *url.URL
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	httpClient2 := *httpClient
	c := &Client{client: &httpClient2}
	c.initialize()
	return c
}

func (c *Client) WithAuthToken(token string) *Client {
	t := c.client.Transport
	if t == nil {
		t = http.DefaultTransport
	}

	// create a new client with the given token
	c.client.Transport = roundTripperFunc(
		func(req *http.Request) (*http.Response, error) {
			req.Header.Set("Authorization", token)
			return t.RoundTrip(req)
		},
	)

	return c
}

func (c *Client) initialize() {
	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse("https://azure.devices.net/")
	}
	c.Configurations = &ConfigurationsService{client: c, baseURL: c.BaseURL}
}

// roundTripperFunc is a custom type that allows us to use a function as an http.RoundTripper
type roundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip is the implementation of the http.RoundTripper interface for roundTripperFunc
func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	// if !strings.HasSuffix(c.BaseURL.Path, "/") {
	// 	return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	// }

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// var buf io.ReadWriter
	// if body != nil {
	// 	buf = &bytes.Buffer{}
	// 	enc := json.NewEncoder(buf)
	// 	enc.SetEscapeHTML(false)
	// 	err := enc.Encode(body)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// print the request headers to the console
	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, fmt.Errorf("error dumping request: %v", err)
	}
	fmt.Println("Request Headers:")
	fmt.Printf("%s\n", dump)

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	fmt.Println("On Do")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fmt.Println("On Do 2")

	// dump the response to the console
	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		return nil, fmt.Errorf("error dumping response: %v", err)
	}
	fmt.Printf("%s\n", dump)

	decErr := json.NewDecoder(res.Body).Decode(v)
	if decErr == io.EOF {
		decErr = nil
	}
	if decErr != nil {
		err = decErr
	}

	return res, nil
}

// func updateManifest(c *configuration.Configuration) error {
// 	// parse manifest file to struct
// 	m := Manifest{}

// 	manifest, err := os.Open(c.Module.Manifest)
// 	if err != nil {
// 		return fmt.Errorf("failed to read manifest: %v", err)
// 	}
// 	defer manifest.Close()

// 	jsonParser := json.NewDecoder(manifest)
// 	if err = jsonParser.Decode(&m); err != nil {
// 		return fmt.Errorf("parsing config file: %v", err)
// 	}

// 	// update the image property to target the new docker image
// 	image := fmt.Sprintf("%s/%s-%s:%s", c.Infra.Registry, c.Image.Repo, c.Device.Arch, c.Id)
// 	if err = setModuleSettings(&m, c.Module.Name, "image", image); err != nil {
// 		return err
// 	}

// 	// update the module version the be 8 char long from uuid
// 	if err = setModuleProperty(&m, c.Module.Name, "version", uuid.New().String()[:8]); err != nil {
// 		return err
// 	}

// 	// update the manifest id
// 	m.Id = fmt.Sprintf("%s-%s", c.Image.Repo, c.Id)

// 	// update the target condition
// 	m.TargetCondition = fmt.Sprintf("tags.application.%s='%s'", c.Module.Name, c.Id)

// 	// write the updated manifest to the file
// 	file, _ := json.MarshalIndent(m, "    ", "    ")
// 	if err = os.WriteFile(c.Module.Manifest, file, 0644); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func setModuleSettings(m *Manifest, moduleName string, property string, value string) error {
// 	if _, ok := m.Content.ModulesContent.EdgeAgent["properties.desired.modules."+moduleName]; !ok {
// 		return fmt.Errorf("module %s is not set as desired module", moduleName)
// 	}
// 	moduleProperties := m.Content.ModulesContent.EdgeAgent["properties.desired.modules."+moduleName].(map[string]interface{})

// 	if _, ok := moduleProperties["settings"]; !ok {
// 		return fmt.Errorf("module manifest %s missing settings field", moduleName)
// 	}
// 	moduleProperties["settings"].(map[string]interface{})[property] = value

// 	return nil
// }

// func setModuleProperty(m *Manifest, moduleName string, property string, value string) error {
// 	if _, ok := m.Content.ModulesContent.EdgeAgent["properties.desired.modules."+moduleName]; !ok {
// 		return fmt.Errorf("module %s is not set as desired moduled", moduleName)
// 	}
// 	moduleProperties := m.Content.ModulesContent.EdgeAgent["properties.desired.modules."+moduleName].(map[string]interface{})

// 	moduleProperties[property] = value

// 	return nil
// }

func prepareRequest(method string, url string, token string) (*http.Request, error) {
	// add the Authorization header to the request
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Authorization", token)
	return r, nil
}

func apiCall(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
