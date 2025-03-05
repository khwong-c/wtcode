package config

import (
	"flag"
	"os"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigdotenv"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/juju/errors"
	"github.com/samber/do"
)

type Env string

const (
	EnvDevelopment Env = "development"
	EnvTest        Env = "test"
	EnvStaging     Env = "staging"
	EnvProduction  Env = "production"
)

type Config struct {
	Env       Env `default:"development" usage:"Environment"`
	HTTPPort  int `default:"8086" usage:"Server port"`
	SQLTarget struct {
		Default string `default:"memory" usage:"Default SQL DSN"`
	}
}

func DefaultLoaderConfig() aconfig.Config {
	return aconfig.Config{
		SkipEnv:   true,
		SkipFlags: false,
		Files:     []string{"config.yaml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
			".env":  aconfigdotenv.New(),
		},
	}
}

func LoadConfig(*do.Injector) (*Config, error) {
	newCfg := Config{}
	loaderConfig := DefaultLoaderConfig()
	loader := aconfig.LoaderFor(&newCfg, loaderConfig)
	if err := loader.Load(); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		return nil, errors.Trace(err)
	}
	return &newCfg, nil
}
