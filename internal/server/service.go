package server

import (
	"context"
	"errors"

	project "github.com/ionut-maxim/insider-project"
	"github.com/ionut-maxim/insider-project/internal/worker"
)

var _ StrictServerInterface = (*service)(nil)

// TODO: Add an ErrorHandlerFunc to respond with proper json messages
type service struct {
	store  project.MessageStore
	worker *worker.Worker
}

func newService(store project.MessageStore, worker *worker.Worker) *service {
	return &service{
		store:  store,
		worker: worker,
	}
}

func (s *service) MessagesSent(ctx context.Context, request MessagesSentRequestObject) (MessagesSentResponseObject, error) {
	var result []Message

	if request.Params.Limit == nil {
		v := 10
		request.Params.Limit = &v
	}
	if request.Params.Offset == nil {
		v := 0
		request.Params.Offset = &v
	}

	messages, err := s.store.Sent(ctx, *request.Params.Limit, *request.Params.Offset)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		id := int(message.ID)
		result = append(result, Message{
			Id:        &id,
			Content:   message.Content,
			CreatedAt: &message.CreatedAt,
			To:        message.To,
			UpdatedAt: &message.UpdatedAt,
		})
	}

	return MessagesSent200JSONResponse(result), nil
}
func (s *service) MessagesAdd(ctx context.Context, request MessagesAddRequestObject) (MessagesAddResponseObject, error) {
	if request.Body == nil {
		return nil, errors.New("body is required")
	}

	if len(request.Body.Content) > 400 {
		return nil, errors.New("content is too large")
	}

	if err := s.store.Add(ctx, project.AddMessageRequest{
		To:      request.Body.To,
		Content: request.Body.Content,
	}); err != nil {
		return nil, err
	}

	return MessagesAdd201Response{}, nil
}
func (s *service) WorkerStart(_ context.Context, _ WorkerStartRequestObject) (WorkerStartResponseObject, error) {
	if err := s.worker.Start(context.Background()); err != nil {
		return WorkerStart200JSONResponse{}, err
	}
	return WorkerStart200JSONResponse{Message: "Worker started"}, nil
}
func (s *service) WorkerStop(_ context.Context, _ WorkerStopRequestObject) (WorkerStopResponseObject, error) {
	if err := s.worker.Close(); err != nil {
		return WorkerStop200JSONResponse{}, err
	}
	return WorkerStop200JSONResponse{Message: "Worker Stopped"}, nil
}
