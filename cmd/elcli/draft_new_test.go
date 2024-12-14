package elcli

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/unbrikd/edge-leap/internal/configuration"
)

// Helper function to reset global variables before each test
func resetGlobals() {
	cfgFile = ""
	force = false
	viper.Reset()
	os.Remove(cfgFile)
}

// TestConfigFileCreationNoExisting tests creating a config file when it doesn't exist
func TestConfigFileCreationNoExisting(t *testing.T) {
	resetGlobals()
	cfgFile = "testdata/new_config.yaml"
	force = false

	handleConfigFileCreation()
	assert.FileExists(t, cfgFile)
}

// TestConfigFileCreationExistingWithForce tests overwriting an existing file with force
func TestConfigFileCreationExistingWithForce(t *testing.T) {
	resetGlobals()
	cfgFile = "testdata/existing_config.yaml"
	force = true

	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")
	viper.Set("id", "initial-id")
	err := viper.WriteConfig()
	if err != nil {
		t.Fatalf("Failed to setup initial file: %v", err)
	}
	defer os.Remove(cfgFile)

	handleConfigFileCreation()
	assert.FileExists(t, cfgFile)

	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	assert.Empty(t, viper.GetString("id"))
}

func TestNewDraft(t *testing.T) {
	resetGlobals()
	cfgFile = "testdata/new_config.yaml"

	newDraft("test-id")
	assert.FileExists(t, cfgFile)

	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	assert.Equal(t, "test-id", viper.GetString("session"))
	assert.Equal(t, configuration.CONFIG_VERSION, viper.GetString("version"))
}

// Setup function to create necessary directories before tests
func TestMain(m *testing.M) {
	// Create testdata directory if it doesn't exist
	os.MkdirAll("testdata", 0755)

	// Run tests
	exitCode := m.Run()

	// Optional: Clean up testdata directory after tests
	os.RemoveAll("testdata")

	os.Exit(exitCode)
}
