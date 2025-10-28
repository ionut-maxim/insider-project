package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	project "github.com/ionut-maxim/insider-project"
	"github.com/ionut-maxim/insider-project/internal/worker"
	"github.com/ionut-maxim/insider-project/openapi"
)

const shutdownTimeout = 1 * time.Second

type Server struct {
	server *http.Server
	svc    *service
	logger *slog.Logger
}

func New(portNumber int, store project.MessageStore, worker *worker.Worker, logger *slog.Logger) *Server {
	srv := &Server{
		svc:    newService(store, worker),
		server: &http.Server{},
		logger: logger,
	}

	mux := http.NewServeMux()

	docsHandler := openapi.NewHandler("/docs", "Notification service")
	strictHandlers := NewStrictHandler(srv.svc, nil)
	apiHandler := HandlerWithOptions(strictHandlers, StdHTTPServerOptions{BaseURL: "/api/v1"})

	// TODO: Add logging middleware...
	mux.Handle("/docs/", docsHandler)
	mux.Handle("/api/v1/", apiHandler)

	srv.server.Handler = mux
	srv.server.Addr = fmt.Sprintf(":%d", portNumber)

	return srv
}

func (s *Server) Start() error {
	go func() {
		s.server.ListenAndServe()
	}()

	return nil
}

func (s *Server) Close() error {
	s.svc.worker.Close()

	// TODO: For some reason connections are not properly cleaned up and it blocks indefinitely. Not a big issues as Postgres does the cleanup in the end.
	//s.svc.store.Close()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
