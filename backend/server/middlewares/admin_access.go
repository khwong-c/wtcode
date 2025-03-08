package middlewares

import (
	"net/http"

	"github.com/khwong-c/wtcode/authentication"
	"github.com/khwong-c/wtcode/config"
)

func RequireAdminAccess(cfg *config.Config, auth authentication.Authenticator) func(http.Handler) http.Handler {
	headerKey := cfg.AdminKey.Header
	middleware := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			adminKey := r.Header.Get(headerKey)
			if !auth.IsAdmin(adminKey) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return middleware
}
