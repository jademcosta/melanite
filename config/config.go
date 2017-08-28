package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const EnvVarKeyImageSource string = "MELANITE_CONF_IMAGE_SOURCE"

const httpProtocol string = "http://"
const httpsProtocol string = "https://"

type Config struct {
	ImageSource string `yaml:"image_source"`
	Port        string `yaml:"port"`
}

func New(rawConfig []byte) (Config, error) {

	configFromEnv := getConfigValuesFromEnv()

	config, err := parseConfigFromString(rawConfig, &Config{})
	if err != nil {
		return config, err
	}

	config = mergeConfigs(config, configFromEnv)

	if err = valid(config); err != nil {
		return config, err
	}

	return config, nil
}

func parseConfigFromString(rawConfig []byte, config *Config) (Config, error) {
	err := yaml.Unmarshal(rawConfig, config)
	return *config, err
}

func getConfigValuesFromEnv() Config {
	config := Config{}

	if imageSource, ok := os.LookupEnv(EnvVarKeyImageSource); ok {
		config.ImageSource = imageSource
	}
	return config
}

func mergeConfigs(lowPriorityConfig, highPriorityConfig Config) Config {
	if highPriorityConfig.ImageSource != "" {
		lowPriorityConfig.ImageSource = highPriorityConfig.ImageSource
	}
	return lowPriorityConfig
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
