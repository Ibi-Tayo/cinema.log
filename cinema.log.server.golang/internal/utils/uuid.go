package utils

import (
	"fmt"

	"github.com/google/uuid"
)

// ParseUUID safely parses a string to UUID and returns an error if invalid
func ParseUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, fmt.Errorf("empty UUID string")
	}
	
	parsed, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %s", s)
	}
	
	return parsed, nil
}

// ValidateUUID checks if a string is a valid UUID
func ValidateUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// GenerateUUID generates a new random UUID
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// IsNilUUID checks if UUID is nil (zero value)
func IsNilUUID(id uuid.UUID) bool {
	return id == uuid.Nil
}
