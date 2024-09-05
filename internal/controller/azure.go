package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/unbrikd/edge-leap/internal/configuration"
)

type Manifest struct {
	Content struct {
		ModulesContent struct {
			EdgeAgent map[string]interface{} `json:"$edgeAgent"`
		} `json:"modulesContent"`
	} `json:"content"`
	Id              string `json:"id"`
	Priority        int    `json:"priority"`
	TargetCondition string `json:"targetCondition"`
}

// azCtl is a controller for Azure IoT Hub.
type azCtl struct {
	// hubName is the name of the Azure IoT Hub.
	hubName string
	// token is the authentication token for the Azure IoT Hub.
	token string
}

func Azure(hubName string, token string) *azCtl {
	return &azCtl{
		hubName: hubName,
		token:   token,
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

func (az *azCtl) CreateLayeredDeployment(c *configuration.Configuration) error {
	deploymentId := fmt.Sprintf("%s-%s", c.Image.Repo, c.Id)

	url := fmt.Sprintf("https://iot-bems-dev-euw-p01.azure-devices.net/configurations/%s?api-version=2021-04-12", deploymentId)

	// Modify the manifest content
	err := updateManifest(c)
	if err != nil {
		return fmt.Errorf("error updating module manifest file: %v", err)
	}

	// read the updated manifest file
	newManifest, err := os.ReadFile(c.Module.Manifest)
	if err != nil {
		log.Fatalf("failed to read new manifest: %v", err)
	}

	// create the request
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", az.token)
	req.Header.Set("Content-Type", "application/json")

	req.Body = io.NopCloser(bytes.NewReader(newManifest))
	req.Body.Close()

	res, err := apiCall(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized, please set a valid token")
	}

	if res.StatusCode == http.StatusConflict {
		return fmt.Errorf("deployment %s already exists", deploymentId)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected hub behaviour: %v", res.Status)
	}

	return nil
}

func (az *azCtl) DeleteLayeredDeployment(c *configuration.Configuration) error {
	deploymentId := fmt.Sprintf("%s-%s", c.Image.Repo, c.Id)

	url := fmt.Sprintf("https://iot-bems-dev-euw-p01.azure-devices.net/configurations/%s?api-version=2021-04-12", deploymentId)

	req, err := prepareRequest(http.MethodDelete, url, az.token)
	if err != nil {
		return err
	}

	res, err := apiCall(req)
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

func updateManifest(c *configuration.Configuration) error {
	// parse manifest file to struct
	m := Manifest{}

	manifest, err := os.Open(c.Module.Manifest)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %v", err)
	}
	defer manifest.Close()

	jsonParser := json.NewDecoder(manifest)
	if err = jsonParser.Decode(&m); err != nil {
		return fmt.Errorf("parsing config file: %v", err)
	}

	// update the image property to target the new docker image
	image := fmt.Sprintf("%s/%s-%s:%s", c.Infra.Registry, c.Image.Repo, c.Device.Arch, c.Id)
	if err = setModuleSettings(&m, c.Module.Name, "image", image); err != nil {
		return err
	}

	// update the module version the be 8 char long from uuid
	if err = setModuleProperty(&m, c.Module.Name, "version", uuid.New().String()[:8]); err != nil {
		return err
	}

	// update the manifest id
	m.Id = fmt.Sprintf("%s-%s", c.Image.Repo, c.Id)

	// update the target condition
	m.TargetCondition = fmt.Sprintf("tags.application.%s='%s'", c.Module.Name, c.Id)

	// write the updated manifest to the file
	file, _ := json.MarshalIndent(m, "    ", "    ")
	if err = os.WriteFile(c.Module.Manifest, file, 0644); err != nil {
		return err
	}

	return nil
}

func setModuleSettings(m *Manifest, moduleName string, property string, value string) error {
	if _, ok := m.Content.ModulesContent.EdgeAgent["properties.desired.modules."+moduleName]; !ok {
		return fmt.Errorf("module %s is not set as desired moduled", moduleName)
	}
	moduleProperties := m.Content.ModulesContent.EdgeAgent["properties.desired.modules."+moduleName].(map[string]interface{})

	if _, ok := moduleProperties["settings"]; !ok {
		return fmt.Errorf("module manifest %s missing settings field", moduleName)
	}
	moduleProperties["settings"].(map[string]interface{})[property] = value

	return nil
}

func setModuleProperty(m *Manifest, moduleName string, property string, value string) error {
	if _, ok := m.Content.ModulesContent.EdgeAgent["properties.desired.modules."+moduleName]; !ok {
		return fmt.Errorf("module %s is not set as desired moduled", moduleName)
	}
	moduleProperties := m.Content.ModulesContent.EdgeAgent["properties.desired.modules."+moduleName].(map[string]interface{})

	moduleProperties[property] = value

	return nil
}

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
	defer res.Body.Close()

	return res, nil
}
