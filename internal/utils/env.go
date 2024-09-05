package utils

import "os"

// GetEnv returns the value of an environment variable or a fallback value if the environment
// variable is not set.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
