package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/server"
	"github.com/khwong-c/wtcode/tests/fixtures"
	"github.com/khwong-c/wtcode/tooling/di"
)

const adminAPIKey = "admin"

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
	di.InvokeOrProvide(ts.injector, func(injector *do.Injector) (*config.Config, error) {
		cfg, err := fixtures.CreateDefaultConfig(false)
		if err != nil {
			return nil, err
		}
		// Patch the Config for testing.
		cfg.AdminKey.Value = adminAPIKey
		cfg.DBSetup.AutoMigrate = true
		return cfg, nil
	})

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
