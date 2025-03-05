package middlewares

import (
	"net/http"

	"github.com/juju/errors"
	"github.com/unrolled/render"

	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/tooling/log"
)

type errorMessage struct {
	Error   errorType `json:"error"`
	Message string    `json:"message"`
}
type detailedErrorMessage struct {
	errorMessage
	Stack string `json:"stack"`
}

type errorType string

var errLogger = log.NewLogger("server.error")

const (
	ErrorTypeNotFound      errorType = "ErrNotFound"
	ErrorTypeNotValid      errorType = "ErrNotValid"
	ErrorTypeInternalError errorType = "ErrInternalError"
)

const (
	// DebugAPIKey is the header key for developer to enter debug api key.
	DebugAPIKey = "X-Debug-Key" // nolint: gosec
)

func handleErr(
	err error,
	cfg *config.Config,
	renderer *render.Render,
	w http.ResponseWriter,
	r *http.Request,
) {
	// handle common errors
	switch {
	case errors.Is(err, errors.NotImplemented):
		_ = renderer.JSON(w, http.StatusNotImplemented, nil)
		return
	case errors.Is(err, errors.Unauthorized):
		_ = renderer.JSON(w, http.StatusUnauthorized, nil)
		return
	case errors.Is(err, errors.Forbidden):
		_ = renderer.JSON(w, http.StatusForbidden, nil)
		return
	case errors.Is(err, errors.Timeout) || errors.Is(err, http.ErrHandlerTimeout):
		_ = renderer.JSON(w, http.StatusRequestTimeout, nil)
		return
	case errors.Is(err, errors.NotFound):
		_ = renderer.JSON(
			w,
			http.StatusNotFound,
			errorMessage{
				Error:   ErrorTypeNotFound,
				Message: err.Error(),
			},
		)
		return
	case errors.Is(err, errors.NotValid):
		_ = renderer.JSON(
			w,
			http.StatusBadRequest,
			errorMessage{
				Error:   ErrorTypeNotValid,
				Message: err.Error(),
			},
		)
		return
	// Other unknown errors
	default:
		handleInternalErr(err, cfg, renderer, w, r)
		return
	}
}

func RequestInDebugMode(r *http.Request, cfg *config.Config) bool {
	debugKey := r.Header.Get(DebugAPIKey)
	return debugKey != "" && debugKey == cfg.DebugKey
}

func handleInternalErr(
	err error,
	cfg *config.Config,
	renderer *render.Render,
	w http.ResponseWriter,
	r *http.Request,
) {
	debugging := false
	switch cfg.Env {
	case config.EnvDevelopment:
		debugging = true
	case config.EnvTest:
		debugging = true
	case config.EnvStaging:
		debugging = RequestInDebugMode(r, cfg)
	case config.EnvProduction:
		debugging = RequestInDebugMode(r, cfg)
	default:
		debugging = RequestInDebugMode(r, cfg)
	}

	// Put the error in the log
	msg, stack := err.Error(), errors.ErrorStack(err)
	errLogger.Error(msg, "stack", stack)

	// Render the error to the client.
	if !debugging {
		_ = renderer.JSON(
			w,
			http.StatusInternalServerError,
			errorMessage{
				Error:   ErrorTypeInternalError,
				Message: "An internal error occurred",
			},
		)
		return
	}
	_ = renderer.JSON(
		w,
		http.StatusInternalServerError,
		detailedErrorMessage{
			errorMessage: errorMessage{
				Error:   ErrorTypeInternalError,
				Message: msg,
			},
			Stack: stack,
		},
	)
}

func PanicRecovery(cfg *config.Config, rend *render.Render) func(http.Handler) http.Handler {
	middleware := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(rvr)
					}
					if err, ok := rvr.(error); ok {
						handleErr(err, cfg, rend, w, r)
					} else {
						errLogger.Crit("U.N.Known Panic Recovered", "panic", rvr)
					}

					if r.Header.Get("Connection") != "Upgrade" {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return middleware
}
