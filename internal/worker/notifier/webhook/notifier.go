package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ionut-maxim/insider-project/internal/worker"
)

type Notifier struct {
	url    string
	client *http.Client
}

func New(url string) *Notifier {
	return &Notifier{
		client: &http.Client{},
		url:    url,
	}
}

func (n *Notifier) Notify(ctx context.Context, notification worker.Notification) (worker.Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(notification); err != nil {
		return worker.Response{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, n.url, &buf)
	if err != nil {
		return worker.Response{}, err
	}

	httpRes, err := n.client.Do(httpReq)
	if err != nil {
		return worker.Response{}, err
	}
	defer httpRes.Body.Close()

	statusOK := httpRes.StatusCode >= 200 && httpRes.StatusCode < 300
	if !statusOK {
		return worker.Response{}, fmt.Errorf("unable to send notification: status code %d", httpRes.StatusCode)
	}

	var response worker.Response
	if err = json.NewDecoder(httpRes.Body).Decode(&response); err != nil {
		return worker.Response{}, errors.New("unable to parse webhook response")
	}

	return response, nil
}
