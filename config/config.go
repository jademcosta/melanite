package config

import (
	"fmt"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const httpProtocol string = "http://"
const httpsProtocol string = "https://"

type Config struct {
	ImageSource string `yaml:"image_source"`
	Port        string `yaml:"port"`
}

func New(rawConfig []byte, overrideImageSource string) (Config, error) {
	config := &Config{}

	err := yaml.Unmarshal(rawConfig, config)
	if err != nil {
		return *config, err
	}

	if overrideImageSource != "" {
		config.ImageSource = overrideImageSource
	}

	if err = valid(*config); err != nil {
		return *config, err
	}

	return *config, nil
}

func valid(config Config) error {
	if !strings.HasPrefix(config.ImageSource, httpProtocol) &&
		!strings.HasPrefix(config.ImageSource, httpsProtocol) {

		return fmt.Errorf("config: image_source should start with %s or %s",
			httpProtocol, httpsProtocol)
	}
	if _, err := strconv.Atoi(config.Port); config.Port != "" && err != nil {
		return fmt.Errorf("config: port should be an integer")
	}
	return nil
}
