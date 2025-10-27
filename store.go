package insider_project

import (
	"context"
	"time"
)

type Status uint8

const (
	StatusUnsent Status = iota
	StatusSent
)

type Message struct {
	ID        uint64    `json:"id"`
	To        string    `json:"to"`
	Content   string    `json:"content"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddMessageRequest struct {
	To      string
	Content string
}

type MessageStore interface {
	// Sent retrieves all sent messages from MessageStore
	Sent(ctx context.Context) ([]Message, error)

	// Unsent retrieves next two messages that should be sent
	Unsent(ctx context.Context) ([2]Message, error)

	// Add adds a new message to the store
	Add(ctx context.Context, req AddMessageRequest) error

	// Update will set the status on a Message
	Update(ctx context.Context, id uint64, status Status) error
}
