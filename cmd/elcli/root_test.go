package elcli

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

// Helper function to reset global variables before each test
func resetGlobals() {
	cfgFile = ""
	force = false
	viper.Reset()
	os.Remove(cfgFile)
}

// Setup function to create necessary directories before tests
func TestMain(m *testing.M) {
	os.MkdirAll("testdata", 0755)
	exitCode := m.Run()
	os.RemoveAll("testdata")
	os.Exit(exitCode)
}
