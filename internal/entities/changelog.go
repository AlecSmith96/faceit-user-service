package entities

import (
	"github.com/google/uuid"
	"time"
)

// ChangelogEntry is a struct that represents a change to a user entity.
// Could be improved by adding a field for an identifier of the person who made the change
type ChangelogEntry struct {
	UserID     uuid.UUID
	CreatedAt  time.Time
	ChangeType string
}
