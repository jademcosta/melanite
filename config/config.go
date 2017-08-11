package config

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const http_protocol string = "http://"
const https_protocol string = "https://"

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
	if strings.HasPrefix(config.ImageSource, http_protocol) ||
		strings.HasPrefix(config.ImageSource, https_protocol) ||
		config.ImageSource == "" {
		return nil
	}
	return fmt.Errorf("config: image_source should start with %s or %s",
		http_protocol, https_protocol)
}
