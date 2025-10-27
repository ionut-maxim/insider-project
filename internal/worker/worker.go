package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"

	project "github.com/ionut-maxim/insider-project"
)

type Worker struct {
	pollInterval time.Duration

	mu      sync.Mutex
	ticker  *time.Ticker
	stop    chan struct{}
	running bool

	logger *slog.Logger

	store    project.MessageStore
	cache    Cache
	notifier Notifier
}

type Option func(*Worker)

func WithInterval(d time.Duration) Option {
	return func(w *Worker) {
		w.pollInterval = d
	}
}

func New(store project.MessageStore, notifier Notifier, cache Cache, logger *slog.Logger, options ...Option) *Worker {
	worker := &Worker{
		pollInterval: 2 * time.Minute,
		logger:       logger.With("service", "worker"),
		store:        store,
		cache:        cache,
		notifier:     notifier,
	}

	for _, option := range options {
		option(worker)
	}

	return worker
}

func (w *Worker) Start(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return ErrAlreadyRunning
	}

	w.ticker = time.NewTicker(w.pollInterval)
	w.stop = make(chan struct{})
	w.running = true

	w.logger.Info("Starting...")

	go w.start(ctx)

	return nil
}

func (w *Worker) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return ErrAlreadyStopped
	}

	close(w.stop)
	w.running = false

	return nil
}

func (w *Worker) start(ctx context.Context) {
	if err := w.processMessages(ctx); err != nil {
		w.logger.Error("Failed to fetch and send messages", "details", err)
	}

	for {
		select {
		case <-w.ticker.C:
			if err := w.processMessages(ctx); err != nil {
				w.logger.Error("Failed to fetch and send messages", "details", err)
			}
		case <-w.stop:
			w.ticker.Stop()
			w.logger.Info("Shutting down...")
			return
		case <-ctx.Done():
			w.ticker.Stop()
			w.logger.Info("Shutting down due to context cancellation...")
			return

		}
	}
}

func (w *Worker) processMessages(ctx context.Context) error {
	messages, err := w.store.Next(ctx)
	if err != nil {
		return ErrUnableToFetchMessages
	}

	for _, message := range messages {
		response, err := w.notifier.Notify(ctx, Notification{To: message.To, Content: message.Content})
		if err != nil {
			return ErrUnableToNotify
		}
		w.logger.Debug("Notification sent", "to", message.To, "content", message.Content, "response", response)

		if err = w.store.Update(ctx, message.ID, project.StatusSent); err != nil {
			return ErrUnableToUpdateStatus
		}

		if err = w.cache.Set(ctx, response.MessageId); err != nil {
			return ErrUnableToUpdateCache
		}
	}

	return nil
}
