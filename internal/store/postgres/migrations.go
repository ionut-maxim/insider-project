package postgres

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func applyMigrations(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger) error {
	if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("creating migrations table: %w", err)
	}

	names, err := fs.Glob(migrationFS, "migrations/*.sql")
	if err != nil {
		return err
	}

	for _, name := range names {
		if err = migrate(ctx, pool, name, logger); err != nil {
			return fmt.Errorf("migration error '%q': %w", name, err)
		}
	}
	return nil
}

func migrate(ctx context.Context, pool *pgxpool.Pool, name string, logger *slog.Logger) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	var n int
	if err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM migrations WHERE name = $1", name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil
	}

	logger.Debug("Applying migration", "name", name)
	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(ctx, string(buf)); err != nil {
		return err
	}

	if _, err = tx.Exec(ctx, "INSERT INTO migrations (name) VALUES ($1)", name); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
