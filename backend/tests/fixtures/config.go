package fixtures

import (
	"flag"
	"os"

	"github.com/cristalhq/aconfig"
	"github.com/juju/errors"

	"github.com/khwong-c/wtcode/config"
)

func CreateDefaultConfig(withEnv bool) (*config.Config, error) {
	newCfg := config.Config{}
	loaderConfig := config.DefaultLoaderConfig()
	loaderConfig.SkipFlags = true
	loaderConfig.SkipEnv = !withEnv
	loader := aconfig.LoaderFor(&newCfg, loaderConfig)
	if err := loader.Load(); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		return nil, errors.Trace(err)
	}
	return &newCfg, nil
}
