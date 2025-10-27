package postgres

import (
	"context"
	"testing"

	project "github.com/ionut-maxim/insider-project"
)

func Test_Store(t *testing.T) {
	ctx := context.Background()

	// TODO: Add postgres test container
	url := "postgres://postgres:postgres@backend.orb.local:5432/backend?sslmode=disable"

	s, err := New(url)
	if err != nil {
		t.Fatal(err)
	}

	for range 2 {
		if err = s.Add(ctx, project.AddMessageRequest{To: "random recipient", Content: "random content"}); err != nil {
			t.Fatal(err)
		}
	}

	next, err := s.Next(ctx)
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
