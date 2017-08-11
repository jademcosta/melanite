package config

import yaml "gopkg.in/yaml.v2"

type Config struct {
	ImageSource string `yaml:"image_source"`
}

func New(rawConfig []byte) (Config, error) {
	config := &Config{}

	err := yaml.Unmarshal(rawConfig, config)
	if err != nil {
		return *config, err
	}
	return *config, nil
}
