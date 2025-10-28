package postgres

import (
	"context"
	"log/slog"
	"os"
	"testing"

	project "github.com/ionut-maxim/insider-project"
)

func Test_Store(t *testing.T) {
	ctx := context.Background()

	//postgresContainer, err := postgres.Run(context.Background(),
	//	"postgres:18-alpine",
	//	postgres.WithDatabase("backend"),
	//	postgres.WithUsername("postgres"),
	//	postgres.WithPassword("postgres"),
	//)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//defer postgresContainer.Terminate(ctx)
	//
	//url, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	//if err != nil {
	//	t.Fatal(err)
	//}
	url := "postgres://postgres:postgres@backend.orb.local/postgres?sslmode=disable"

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	s, err := New(url, logger)
	if err != nil {
		t.Fatal(err)
	}

	for range 2 {
		if err = s.Add(ctx, project.AddMessageRequest{To: "random recipient", Content: "random content"}); err != nil {
			t.Fatal(err)
		}
	}

	next, err := s.Next(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(next) != 2 {
		t.Fatalf("got %d items, want 2", len(next))
	}

	for _, msg := range next {
		if err = s.Update(ctx, msg.ID, project.StatusSent); err != nil {
			t.Fatal(err)
		}
	}

	messages, err := s.Sent(ctx, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) == 0 {
		t.Fatal("no messages received")
	}
}
