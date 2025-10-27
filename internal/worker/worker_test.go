package worker_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"testing/synctest"
	"time"

	project "github.com/ionut-maxim/insider-project"
	memorystore "github.com/ionut-maxim/insider-project/store/memory"
	"github.com/ionut-maxim/insider-project/worker"
	memorycache "github.com/ionut-maxim/insider-project/worker/cache/memory"
	"github.com/ionut-maxim/insider-project/worker/notifier/stub"
)

func Test_Worker(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx := context.Background()

		store := memorystore.New()
		notifier := &stub.Notifier{}
		cache := memorycache.New()
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

		w := worker.New(store, notifier, cache, logger)

		if err := w.Start(ctx); err != nil {
			t.Fatal(err)
		}

		if err := store.Add(ctx, project.AddMessageRequest{
			To:      "Test Recipient",
			Content: "Test Message",
		}); err != nil {
			t.Fatal(err)
		}

		time.Sleep(5 * time.Minute)

		if err := store.Add(ctx, project.AddMessageRequest{
			To:      "Another Test Recipient",
			Content: "Test after 5 minutes",
		}); err != nil {
			t.Fatal(err)
		}

		time.Sleep(5 * time.Minute)

		if err := w.Stop(); err != nil {
			t.Fatal(err)
		}
	})
}
