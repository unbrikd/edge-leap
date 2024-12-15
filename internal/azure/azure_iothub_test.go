package azure_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unbrikd/edge-leap/internal/azure"
)

func TestSetContent(t *testing.T) {
	c := azure.Configuration{}
	moduleName := "myModule"
	image := "img"
	createOpts := "opts"
	startupOrder := 10
	envVars := map[string]string{
		"ENV_VAR_1": "KEY_VAR_1",
		"ENV_VAR_2": "",
	}

	c.SetContent(moduleName, image, createOpts, startupOrder, envVars)

	modulesContent := c.Content["modulesContent"]
	if _, ok := modulesContent.(map[string]interface{})["$edgeAgent"]; !ok {
		t.Fatal("configuration contents is missing '$edgeAgent' key")
	}

	expectedModulePropertiesKey := fmt.Sprintf("properties.desired.modules.%s", moduleName)
	edgeAgent := modulesContent.(map[string]interface{})["$edgeAgent"]
	if _, ok := edgeAgent.(map[string]interface{})[expectedModulePropertiesKey]; !ok {
		t.Fatalf("configuration contents is missing '%s' key", expectedModulePropertiesKey)
		return
	}
	moduleProperties := edgeAgent.(map[string]interface{})[expectedModulePropertiesKey]
	assert.Contains(t, moduleProperties, "settings")
	assert.Contains(t, moduleProperties, "startupOrder")
	assert.Contains(t, moduleProperties, "env")

	settings := moduleProperties.(map[string]interface{})["settings"]
	assert.Contains(t, settings, "image")
	assert.Contains(t, settings, "createOptions")

	env := moduleProperties.(map[string]interface{})["env"]
	for k, v := range envVars {
		assert.Contains(t, env, k)

		gotV := env.(map[string]interface{})[k].(struct {
			Value string `json:"value"`
		}).Value
		assert.Equal(t, v, gotV)
	}

	assert.Equal(t, moduleProperties.(map[string]interface{})["type"], "docker")
	assert.Equal(t, moduleProperties.(map[string]interface{})["status"], "running")
	assert.Equal(t, moduleProperties.(map[string]interface{})["restartPolicy"], "always")
	assert.Equal(t, moduleProperties.(map[string]interface{})["version"], "1.0")
}
