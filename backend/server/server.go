package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/inconshreveable/log15"
	"github.com/juju/errors"
	"github.com/samber/do"

	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/tooling/di"
	"github.com/khwong-c/wtcode/tooling/log"
)

type Server struct {
	http.Server
	Injector *do.Injector
	Handler  http.Handler

	config *config.Config
	logger log15.Logger
}

func (s *Server) Serve() {
	go func() {
		err := s.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Server error", "err", err, "stack", errors.ErrorStack(err))
		}
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.Server.Shutdown(ctx)

	if err != nil {
		s.logger.Error("Server shutdown error", "err", err)
		return errors.Trace(err)
	}
	return nil
}

func (s *Server) WaitForSignal() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

func CreateServer(injector *do.Injector) (*Server, error) {
	var err error
	server := &Server{
		Injector: injector,
		config:   di.InvokeOrProvide(injector, config.LoadConfig),
		logger:   log.NewLogger("server"),
	}
	server.Addr = fmt.Sprintf(":%d", server.config.HTTPPort)

	if err = server.createStack(injector); err != nil {
		return nil, errors.Trace(err)
	}
	if server.Handler, err = server.createRoute(injector); err != nil {
		return nil, errors.Trace(err)
	}
	return server, nil
}

func (s *Server) createStack(injector *do.Injector) error {
	return errors.NotImplemented
}

func (s *Server) createRoute(injector *do.Injector) (http.Handler, error) {
	r := chi.NewMux()
	return r, errors.NotImplemented
}
