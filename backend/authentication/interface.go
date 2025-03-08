package authentication

import (
	"github.com/samber/do"

	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/tooling/di"
)

type Authenticator interface {
	IsAdmin(string) bool
}

type apiKeyAuthenticator struct {
	injector *do.Injector
	adminKey string
}

func CreateAPIKeyAuthenticator(injector *do.Injector) (Authenticator, error) {
	cfg := di.InvokeOrProvide(injector, config.LoadConfig)
	adminKey := cfg.AdminKey.Value
	return &apiKeyAuthenticator{
		injector: injector,
		adminKey: adminKey,
	}, nil
}
