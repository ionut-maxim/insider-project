package worker

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Cache interface {
	// Set will add or update the cache with the current time for the message
	Set(ctx context.Context, id uuid.UUID) error

	// Get retrieves the sentAt time for the specified message ID
	Get(ctx context.Context, id uuid.UUID) (time.Time, error)
}
