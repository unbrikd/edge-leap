package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// A Client manages communication with the GitHub API.
type Client struct {
	clientMu sync.Mutex   // clientMu protects the client during calls that modify the CheckRedirect func.
	client   *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to the public GitHub API, but can be
	// set to a domain endpoint to use with GitHub Enterprise. BaseURL should
	// always be specified with a trailing slash.
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
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
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
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decErr := json.NewDecoder(res.Body).Decode(v)
	if decErr == io.EOF {
		decErr = nil
	}
	if decErr != nil {
		err = decErr
		return nil, err
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
