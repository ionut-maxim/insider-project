package insider_project

import (
	"context"
	"io"
	"time"
)

type Status uint8

const (
	StatusUnsent Status = iota
	StatusSent
)

type Message struct {
	ID      uint32 `json:"id"`
	To      string `json:"to"`
	Content string `json:"content"`
	Status  Status `json:"status"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AddMessageRequest struct {
	To      string
	Content string
}

type MessageStore interface {
	// Sent retrieves sent messages from the store
	Sent(ctx context.Context, limit, offset int) ([]Message, error)

	// Next retrieves the next 2 messages to be sent
	Next(ctx context.Context, limit int) ([]Message, error)

	// Add adds a new message to the store
	Add(ctx context.Context, req AddMessageRequest) error

	// Update will set the status on a Message
	Update(ctx context.Context, id uint32, status Status) error

	io.Closer
}
