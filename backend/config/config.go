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
	Env       Env    `default:"development" usage:"Environment"`
	DebugKey  string `usage:"API key to enable debug mode"`
	HTTPPort  int    `default:"8086" usage:"Server port"`
	SQLTarget struct {
		Default string `default:"sqlite3::memory:" usage:"Default SQL DSN"`
	}
	DBSetup struct {
		AutoMigrate bool `default:"false" usage:"Auto migrate database"`
	}
	AdminKey struct {
		Header string `default:"X-Admin-Key" usage:"Header field for the Admin Mode API key"`
		Value  string `usage:"Admin Mode API Key"`
	} `usage:"Admin Mode API key setup"`
}

func DefaultLoaderConfig() aconfig.Config {
	return aconfig.Config{
		SkipEnv:   true,
		SkipFlags: false,
		FileFlag:  "config",
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
