package memory_test

import (
	"context"
	"fmt"
	"testing"

	project "github.com/ionut-maxim/insider-project"
	"github.com/ionut-maxim/insider-project/store/memory"
)

func Test_MemoryStore(t *testing.T) {
	s := memory.New()
	ctx := context.TODO()

	for i := range 10 {
		req := project.AddMessageRequest{
			To:      fmt.Sprintf("Number %d", i+1),
			Content: fmt.Sprintf("Content %d", i+1),
		}
		if err := s.Add(ctx, req); err != nil {
			t.Error(err)
		}
	}

	for i := range 5 {
		if err := s.Update(ctx, uint64(i+1), project.StatusSent); err != nil {
			t.Error(err)
		}
	}

	sent, err := s.Sent(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(sent) != 5 {
		t.Errorf("Sent count mismatch: got %d, want %d", len(sent), 5)
	}

	unsent, err := s.Unsent(ctx)
	if err != nil {
		t.Error(err)
	}
	if unsent[0].ID != 6 {
		t.Errorf("Unsent id mismatch: got %d, want %d", unsent[0].ID, 6)
	}
}
