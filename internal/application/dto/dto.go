package dto

import "github.com/google/uuid"

// ParseUUID parses a UUID string and returns a uuid.UUID object.
func ParseUUID(uuidStr string) (uuid.UUID, error) {
	return uuid.Parse(uuidStr)
}
