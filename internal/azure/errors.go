package azure

import "fmt"

type ConfigExistsError struct {
	Id string
}

// Implement the error interface
func (e *ConfigExistsError) Error() string {
	return fmt.Sprintf("configuration '%s' already exists", e.Id)
}
