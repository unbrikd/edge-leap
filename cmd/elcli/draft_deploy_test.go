package elcli

import (
	"net/url"

	"github.com/unbrikd/edge-leap/internal/azure"
)

type MockClient struct {
	AuthToken string
	BaseURL   *url.URL
}

func (m *MockClient) WithAuthToken(token string) *azure.Client {
	m.AuthToken = token
	return &azure.Client{BaseURL: m.BaseURL}
}

func (m *MockClient) SetBaseURL(baseURL string) {
	m.BaseURL, _ = url.Parse(baseURL)
}

type MockReleaser struct {
	Client        *azure.Client
	DeviceName    string
	ModuleName    string
	Configuration azure.Configuration
}

func (m *MockReleaser) SetModuleOnDevice(deviceName, moduleName, id string) error {
	m.DeviceName = deviceName
	m.ModuleName = moduleName
	return nil
}

func (m *MockReleaser) ReleaseModule(config *azure.Configuration) error {
	m.Configuration = *config
	return nil
}
