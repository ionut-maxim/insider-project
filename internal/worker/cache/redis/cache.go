package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

type Options func(*Cache)

func WithTTL(ttl time.Duration) Options {
	return func(c *Cache) {
		c.ttl = ttl
	}
}

func New(url string, db int, options ...Options) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: url,
		DB:   db,
	})

	c := &Cache{
		client: client,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *Cache) Set(ctx context.Context, id uuid.UUID) error {
	return c.client.Set(ctx, id.String(), time.Now(), c.ttl).Err()
}

func (c *Cache) Get(ctx context.Context, id uuid.UUID) (time.Time, error) {
	val, err := c.client.Get(ctx, id.String()).Result()
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, val)
}
