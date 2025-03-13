package tests

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cristalhq/aconfig"
	"github.com/juju/errors"
	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/server"
	"github.com/khwong-c/wtcode/tooling/di"
)

const adminAPIKey = "admin"

func skipConfigFlags(*do.Injector) (*config.Config, error) {
	newCfg := config.Config{}
	loaderConfig := config.DefaultLoaderConfig()
	loaderConfig.SkipFlags = true
	loader := aconfig.LoaderFor(&newCfg, loaderConfig)
	if err := loader.Load(); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		return nil, errors.Trace(err)
	}
	return patchCfg(&newCfg), nil
}

// Patch the config for testing.
func patchCfg(cfg *config.Config) *config.Config {
	cfg.AdminKey.Value = adminAPIKey
	return cfg
}

type SmokeTestSuite struct {
	suite.Suite
	injector *do.Injector
	svr      *server.Server
}

func TestSmokeTests(t *testing.T) {
	suite.Run(t, new(SmokeTestSuite))
}

func (ts *SmokeTestSuite) SetupSuite() {
	ts.injector = di.CreateInjector(false, false)
	di.InvokeOrProvide(ts.injector, skipConfigFlags)

	if !ts.NotPanics(func() {
		ts.svr = di.InvokeOrProvide(ts.injector, server.CreateServer)
	}) {
		ts.FailNow("Failed to create server")
	}
}

func (ts *SmokeTestSuite) TearDownSuite() {
	ts.NoError(ts.injector.Shutdown())
}

func (ts *SmokeTestSuite) TestHealthEndpoint() {
	ts.HTTPBodyContains(
		ts.svr.Handler.ServeHTTP,
		http.MethodGet,
		"/health",
		nil,
		".",
	)
}

func (ts *SmokeTestSuite) TestAdminProtect() {
	// Error if no admin key is provided.
	ts.HTTPError(ts.svr.Handler.ServeHTTP, http.MethodGet, "/is_admin", nil, http.StatusUnauthorized)

	// Success if admin key is provided.
	cfg := di.Invoke[*config.Config](ts.injector)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/is_admin", nil)
	req.Header.Add(
		cfg.AdminKey.Header, cfg.AdminKey.Value,
	)
	ts.svr.Handler.ServeHTTP(w, req)
	ts.Equal(http.StatusOK, w.Code)
}
