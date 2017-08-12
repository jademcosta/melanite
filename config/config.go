package config

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const httpProtocol string = "http://"
const httpsProtocol string = "https://"

type Config struct {
	ImageSource string `yaml:"image_source"`
}

func New(rawConfig []byte) (Config, error) {
	config := &Config{}

	err := yaml.Unmarshal(rawConfig, config)
	if err != nil {
		return *config, err
	}

	if err = valid(*config); err != nil {
		return *config, err
	}

	return *config, nil
}

func valid(config Config) error {
	if strings.HasPrefix(config.ImageSource, httpProtocol) ||
		strings.HasPrefix(config.ImageSource, httpsProtocol) {
		return nil
	}
	invalidImageSourceErrorMessage := "config: image_source should start with %s or %s"
	return fmt.Errorf(invalidImageSourceErrorMessage,
		httpProtocol, httpsProtocol)
}
