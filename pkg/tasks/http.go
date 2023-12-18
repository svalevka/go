package tasks

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// HTTPServer is a task that runs an HTTP(S) server.
type HTTPServer struct {
	// Name identifies this Task in the TaskStatus calls from the Task Runner.
	Name string

	// Addr is the host:port where this HTTP server will listen.
	Addr string

	// Handler is the HTTP callback used to serve new requests.
	Handler http.Handler
}

func (h *HTTPServer) TaskName() string {
	return "HTTPServer(" + h.Name + "," + h.Addr + ")"
}

func (h *HTTPServer) RunTask(ctx context.Context) error {
	s := &http.Server{
		Addr:    h.Addr,
		Handler: h.Handler,
	}

	go func() {
		<-ctx.Done()
		s.Shutdown(ctx)
	}()

	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return ctx.Err()
	} else if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	return nil
}
