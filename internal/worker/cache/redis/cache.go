package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Cache struct {
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Set(_ context.Context, id uuid.UUID) error {
	return nil
}

func (c *Cache) Get(_ context.Context, id uuid.UUID) (time.Time, error) {
	return time.Time{}, nil
}
