package memory

import (
	"context"
	"sync"
	"time"

	project "github.com/ionut-maxim/insider-project"
)

type Store struct {
	mu       sync.RWMutex
	messages []project.Message
}

func New() *Store {
	return &Store{
		messages: []project.Message{},
	}
}

func (s *Store) Sent(_ context.Context) ([]project.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var messages []project.Message
	for _, m := range s.messages {
		if m.Status != project.StatusSent {
			continue
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (s *Store) Unsent(_ context.Context) ([]project.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var messages []project.Message
	count := 0

	for _, m := range s.messages {
		if m.Status != project.StatusUnsent {
			continue
		}

		messages = append(messages, m)
		count++

		if count == 2 {
			break
		}

	}

	return messages, nil
}

func (s *Store) Add(_ context.Context, req project.AddMessageRequest) error {
	s.mu.Lock()
	id := uint64(len(s.messages) + 1)
	message := project.Message{
		ID:        id,
		To:        req.To,
		Content:   req.Content,
		Status:    project.StatusUnsent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.messages = append(s.messages, message)
	s.mu.Unlock()

	return nil
}

func (s *Store) Update(_ context.Context, id uint64, status project.Status) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.messages {
		if s.messages[i].ID == id {
			s.messages[i].Status = status
			s.messages[i].UpdatedAt = time.Now()
		}
	}

	return nil
}
