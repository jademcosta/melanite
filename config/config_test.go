package config_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jademcosta/melanite/config"
	"github.com/stretchr/testify/assert"
)

const configFixtureFilesFolder = "../test/config_files/"

func TestConfigInitializer(t *testing.T) {
	configContent, err :=
		ioutil.ReadFile(fmt.Sprintf("%sfull_correct_config.yaml", configFixtureFilesFolder))
	if err != nil {
		panic(err)
	}

	configuration, err := config.New(configContent)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, "http://example.com",
		configuration.ImageSource, "The image sources should be equal")
}
