package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/inconshreveable/log15"
	"github.com/juju/errors"
	"github.com/samber/do"
	"github.com/unrolled/render"

	"github.com/khwong-c/wtcode/authentication"
	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/server/middlewares"
	"github.com/khwong-c/wtcode/tooling/di"
	"github.com/khwong-c/wtcode/tooling/log"
)

type Server struct {
	http.Server
	injector *do.Injector
	Handler  http.Handler

	config *config.Config
	logger log15.Logger
	render *render.Render

	auth authentication.Authenticator
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
	const shutdownDuration = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownDuration)
	defer cancel()
	err := s.Server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Server shutdown error", "err", err)
		return errors.Trace(err)
	}
	return nil
}

func CreateServer(injector *do.Injector) (*Server, error) {
	var err error
	server := &Server{
		injector: injector,
		config:   di.InvokeOrProvide(injector, config.LoadConfig),
		logger:   log.NewLogger("server"),
		render:   render.New(),
		auth:     di.InvokeOrProvide(injector, authentication.CreateAPIKeyAuthenticator),
	}
	server.Addr = fmt.Sprintf(":%d", server.config.HTTPPort)

	if server.Handler, err = server.createRoute(); err != nil {
		return nil, errors.Trace(err)
	}
	return server, nil
}

func (s *Server) createRoute() (http.Handler, error) { //nolint:unparam
	// TODO: How to specify the server we want? Is it DI / Compile time config?
	r := chi.NewMux()
	r.Use(middlewares.PanicRecovery(s.config, s.render))
	r.Use(chiMiddleware.Heartbeat("/health"))
	r.With(middlewares.RequireAdminAccess(s.config, s.auth)).Get(
		"/is_admin",
		func(w http.ResponseWriter, _ *http.Request) {
			if err := s.render.JSON(w, http.StatusOK, map[string]interface{}{
				"is_admin": true,
			}); err != nil {
				panic(err)
			}
		},
	)
	return r, nil
}
