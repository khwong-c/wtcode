package tests

import (
	"flag"
	"os"
	"testing"

	"github.com/cristalhq/aconfig"
	"github.com/juju/errors"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"

	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/server"
	"github.com/khwong-c/wtcode/tooling/di"
)

func SkipConfigFlags(*do.Injector) (*config.Config, error) {
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
	return &newCfg, nil
}

func TestBootServer(t *testing.T) {
	injector := di.CreateInjector(false, false)
	di.InvokeOrProvide(injector, SkipConfigFlags)
	assert.NotPanics(t, func() {
		di.InvokeOrProvide(injector, server.CreateServer)
	})
}
