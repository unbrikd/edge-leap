package utils

import (
	"fmt"
	"os"
	"strings"
)

// GetEnv returns the value of an environment variable or a fallback value if the environment
// variable is not set.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// StringArraySplitToMap converts an array of strings into a map of strings using a separator.
func StringArraySplitToMap(arr []string, sep string) (map[string]string, error) {
	m := make(map[string]string)
	for _, v := range arr {
		p := strings.Split(v, sep)
		if len(p) != 2 {
			return nil, fmt.Errorf("invalid key-value pair: %s", v)
		}

		m[p[0]] = p[1]
	}

	return m, nil
}
