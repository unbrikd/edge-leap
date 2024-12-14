package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringArraySplitToMap(t *testing.T) {
	expectedMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	arr := []string{"key1=value1", "key2=value2"}
	m, err := StringArraySplitToMap(arr, "=")

	assert.Nil(t, err)
	assert.Equal(t, expectedMap, m)
}
