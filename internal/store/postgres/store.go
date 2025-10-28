package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	project "github.com/ionut-maxim/insider-project"
)

type Store struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func New(uri string, logger *slog.Logger) (*Store, error) {
	dbCtx := context.Background()
	pool, err := pgxpool.New(dbCtx, uri)
	if err != nil {
		return nil, err
	}

	mLogger := logger.With("service", "migrations")
	if err = applyMigrations(dbCtx, pool, mLogger); err != nil {
		return nil, errors.New("failed to apply postgres migrations")
	}

	return &Store{pool: pool, logger: logger.With("service", "store")}, nil
}

func (s *Store) Sent(ctx context.Context, limit, offset int) ([]project.Message, error) {
	query := `
        SELECT id, to_recipient, content, status, created_at, updated_at
        FROM messages
        WHERE status = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := s.pool.Query(ctx, query, project.StatusSent, limit, offset)
	if err != nil {
		return nil, err
	}

	var messages []project.Message
	for rows.Next() {
		var m project.Message
		if err = rows.Scan(&m.ID, &m.To, &m.Content, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (s *Store) Next(ctx context.Context, limit int) ([]project.Message, error) {
	// Little comment here for running multiple workers:
	// FOR UPDATE locks the two rows we're querying - fixes concurrency issues
	// SKIP LOCKED continues over the skipped rows and selects the next 2 in the result set
	query := `
			SELECT id, to_recipient, content, status, created_at, updated_at
			FROM messages
			WHERE status = 0
			ORDER BY created_at
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		`

	rows, err := s.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	var messages []project.Message
	for rows.Next() {
		var m project.Message
		if err = rows.Scan(&m.ID, &m.To, &m.Content, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (s *Store) Add(ctx context.Context, req project.AddMessageRequest) error {
	query := `
        INSERT INTO messages (to_recipient, content, status)
        VALUES ($1, $2, $3)
    `

	_, err := s.pool.Exec(ctx, query, req.To, req.Content, project.StatusUnsent)
	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, id uint32, status project.Status) error {
	query := `
        UPDATE messages
        SET status = $1,
            updated_at = $2
        WHERE id = $3
    `

	now := time.Now()
	_, err := s.pool.Exec(ctx, query, status, now, id)
	if err != nil {
		return err
	}

	return nil
}
