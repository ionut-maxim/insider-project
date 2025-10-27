package worker

import (
	"context"

	"github.com/google/uuid"
)

type Notification struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type Response struct {
	Message   string    `json:"message"`
	MessageId uuid.UUID `json:"messageId"`
}

type Notifier interface {
	Notify(ctx context.Context, notification Notification) (Response, error)
}
