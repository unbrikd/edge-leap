package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// A Client manages communication with the Azure API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client
	// Base URL for API requests.
	BaseURL *url.URL
	// Configurations service for the Azure IoT Hub API.
	Configurations *ConfigurationsService
}

type Response struct {
	Response *http.Response
}

// service is a genetic struct that abstracts the interaction with Azure resources API.
type service struct {
	client  *Client
	BaseURL *url.URL
}

// NewClient returns a new Azure API client. If a nil httpClient is provided a new http.Client will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	httpClient2 := *httpClient
	c := &Client{client: &httpClient2}
	c.initialize()
	return c
}

// initialize initializes the client with the default values.
func (c *Client) initialize() {
	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse("https://azure.devices.net/")
	}
	c.Configurations = &ConfigurationsService{client: c, BaseURL: c.BaseURL}
}

// WithAuthToken sets the Authorization header for the client.
func (c *Client) WithAuthToken(token string) *Client {
	t := c.client.Transport
	if t == nil {
		t = http.DefaultTransport
	}

	c.client.Transport = roundTripperFunc(
		func(req *http.Request) (*http.Response, error) {
			req.Header.Set("Authorization", token)
			return t.RoundTrip(req)
		},
	)

	return c
}

// roundTripperFunc is a custom type that allows the use a function as an http.RoundTripper
type roundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip is the implementation of the http.RoundTripper interface for roundTripperFunc
func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, in which case it is resolved relative
// to the BaseURL of the Client. The body parameter can be of any type. If the body is not nil, it will be JSON encoded
// and included in the request.
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

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v is nil, the API response is discarded.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if v != nil {
		err = json.NewDecoder(res.Body).Decode(v)
		if err != nil && err != io.EOF {
			return nil, err
		}
	}

	return res, nil
}

// Expect checks if the response status code is in the list of expected status codes. If the status code is not in the
// list, an error is returned with the response status message.
func (r *Response) Expect(statusCode ...int) error {
	// check if response status code is in the list of expected status codes
	for _, code := range statusCode {
		if r.Response.StatusCode == code {
			return nil
		}
	}

	return fmt.Errorf("%s", r.Response.Status)
}

// Is checks if the response status code is equal to the provided status code.
func (r *Response) Is(statusCode int) bool {
	return r.Response.StatusCode == statusCode
}
