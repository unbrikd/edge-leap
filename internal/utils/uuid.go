package utils

import "github.com/google/uuid"

func ShortUuid() string {
	id := uuid.New().String()
	return id[:8]
}
