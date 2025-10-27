package stub

import (
	"context"

	"github.com/google/uuid"

	"github.com/ionut-maxim/insider-project/worker"
)

type Notifier struct {
}

func (n *Notifier) Notify(_ context.Context, _ worker.Notification) (worker.Response, error) {
	return worker.Response{
		Message:   "Accepted",
		MessageId: uuid.New(),
	}, nil
}
