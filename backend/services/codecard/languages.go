package codecard

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed data/CodeSupported.yaml
var supportedLanguagesRaw []byte

type SupportedLanguages map[string][]string

func getSupportedLanguages() (SupportedLanguages, error) {
	var sl SupportedLanguages
	err := yaml.Unmarshal(supportedLanguagesRaw, &sl)
	if err != nil {
		return nil, err
	}
	return sl, nil
}
