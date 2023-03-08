package utils

import "github.com/google/uuid"

// GenerateUUID generates a psuedo random uuid
func GenerateUUID() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
