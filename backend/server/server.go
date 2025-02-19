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
	"github.com/juju/errors"
	"github.com/samber/do"

	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/tooling/di"
)

type Server struct {
	http.Server
	Injector *do.Injector
	Handler  http.Handler

	Config *config.Config
}

func (s *Server) Serve() {
	go func() {
		err := s.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			// TODO: Log error
		}
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.Server.Shutdown(ctx)

	if err != nil {
		return err
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
		Config:   di.InvokeOrProvide(injector, config.LoadConfig),
	}
	server.Addr = fmt.Sprintf(":%d", server.Config.HTTPPort)

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
