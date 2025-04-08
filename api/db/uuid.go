package db

import (
	"github.com/segmentio/ksuid"
)

type UUID struct {
	ksuid.KSUID
}

// Generates a new, wrapped KSUID. In the strange case that random bytes can't be read, it will panic.
func NewUUID() UUID {
	return UUID{KSUID: ksuid.New()}
}

// ParseUUID parses a UUID from a string. If the string is not a valid UUID, it will return an error.
func ParseUUID(s string) (UUID, error) {
	ksuid, err := ksuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	return UUID{KSUID: ksuid}, nil
}
