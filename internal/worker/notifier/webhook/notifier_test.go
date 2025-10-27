package webhook_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/ionut-maxim/insider-project/worker"
	"github.com/ionut-maxim/insider-project/worker/notifier/webhook"
)

func Test_Notifier(t *testing.T) {
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Test message", "messageId": "42049274-e9c0-46cf-b341-7a2022f6dba4"}`))
	}))

	n := webhook.New(srv.URL)

	res, err := n.Notify(ctx, worker.Notification{
		To:      "Test",
		Content: "Test content",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.MessageId == uuid.Nil {
		t.Fatal("expected message id")
	}
}
