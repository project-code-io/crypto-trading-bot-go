package generator

import (
	"fmt"

	"github.com/google/uuid"
)

// RandomUUIDGenerator generates a random ID based on a UUID.
type RandomUUIDGenerator struct{}

// GenerateID will generate a random ID with the prefix + ":" added before.
func (r *RandomUUIDGenerator) GenerateID(prefix string) string {
	randomID := uuid.New()

	return fmt.Sprintf("%s:%s", prefix, &randomID)
}
