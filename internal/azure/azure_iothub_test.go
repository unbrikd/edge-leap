package azure_test

import (
	"fmt"
	"testing"

	"github.com/unbrikd/edge-leap/internal/azure"
)

func TestSetContent(t *testing.T) {
	c := azure.Configuration{}
	moduleName := "myModule"
	image := "img"
	createOpts := "opts"
	startupOrder := "so"
	envVars := map[string]string{
		"ENV_VAR_1": "KEY_VAR_1",
		"ENV_VAR_2": "",
	}

	c.SetContent(moduleName, image, createOpts, startupOrder, envVars)

	// .(map[string]interface{})["$edgeAgent"].(map[string]interface{})["properties.desired.mod"].(map[string]interface{})
	modulesContent := c.Content["modulesContent"]
	edgeAgent := modulesContent.(map[string]interface{})["$edgeAgent"]

	expectedModulePropertiesKey := fmt.Sprintf("properties.desired.modules.%s", moduleName)
	if _, ok := edgeAgent.(map[string]interface{})[expectedModulePropertiesKey]; !ok {
		t.Fatalf("configuration contents is missing '%s' key", expectedModulePropertiesKey)
		return
	}

	moduleProperties := edgeAgent.(map[string]interface{})[expectedModulePropertiesKey]
	if _, ok := moduleProperties.(map[string]interface{})["settings"]; !ok {
		t.Fatal("configuration module properties is missing 'settings' key")
	}

	if _, ok := moduleProperties.(map[string]interface{})["startupOrder"]; !ok {
		t.Fatal("configuration module properties is missing 'startupOrder' key")
	}

	if _, ok := moduleProperties.(map[string]interface{})["env"]; !ok {
		t.Fatal("configuration module properties is missing 'env' key")
	}

	settings := moduleProperties.(map[string]interface{})["settings"]
	if _, ok := settings.(map[string]string)["image"]; !ok {
		t.Fatal("configuration settings sey is missing 'image' key")
	}

	if _, ok := settings.(map[string]string)["createOptions"]; !ok {
		t.Fatal("configuration settings is missing 'createOptions' key")
	}

	env := moduleProperties.(map[string]interface{})["env"]
	for k, v := range envVars {
		if _, ok := env.(map[string]interface{})[k]; !ok {
			t.Fatalf("configuration env is missing '%s' key", k)
		}

		gotV := env.(map[string]interface{})[k].(struct {
			Value string `json:"value"`
		}).Value
		if gotV != v {
			t.Errorf("expected '%s'='%s' got '%s'= '%s'", k, v, k, gotV)
		}
	}

}
