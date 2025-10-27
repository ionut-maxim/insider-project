package memory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Cache struct {
	mutex     sync.RWMutex
	responses map[uuid.UUID]time.Time
}

func New() *Cache {
	return &Cache{
		responses: make(map[uuid.UUID]time.Time),
	}
}

func (c *Cache) Set(_ context.Context, id uuid.UUID) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.responses[id] = time.Now()
	return nil
}

func (c *Cache) Get(_ context.Context, id uuid.UUID) (time.Time, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.responses[id], nil
}
