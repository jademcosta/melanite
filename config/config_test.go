package config_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jademcosta/melanite/config"
	"github.com/stretchr/testify/assert"
)

const configFixtureFilesFolder = "../test/config_files/"

func TestConfigInitializerForValidImageSourceKey(t *testing.T) {

	var imageSourceOnConfigTests = []struct {
		filename            string
		expectedImageSource string
		testMessage         string
	}{
		{"full_correct_config.yaml", "http://example.com",
			"the image source should be the equal"},
	}

	for _, testCase := range imageSourceOnConfigTests {
		configContent, err :=
			ioutil.ReadFile(fmt.Sprintf("%s%s", configFixtureFilesFolder, testCase.filename))
		if err != nil {
			panic(err)
		}

		configuration, err := config.New(configContent)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, testCase.expectedImageSource,
			configuration.ImageSource, testCase.testMessage)
	}
}

func TestConfigInitializerForInvalidImageSourceKey(t *testing.T) {

	var imageSourceOnConfigTests = []struct {
		filename             string
		expectedErrorMessage string
		testMessage          string
	}{
		{"incorrect_image_source_config.yaml",
			"config: image_source should start with http:// or https://",
			"Should return an error"},
		{"empty_config.yaml",
			"config: image_source should start with http:// or https://",
			"Should return an error"},
	}

	for _, testCase := range imageSourceOnConfigTests {
		configContent, err :=
			ioutil.ReadFile(fmt.Sprintf("%s%s", configFixtureFilesFolder, testCase.filename))
		if err != nil {
			panic(err)
		}

		_, err = config.New(configContent)

		assert.Error(t, err, testCase.testMessage)
		assert.Equal(t, err.Error(),
			testCase.expectedErrorMessage, testCase.testMessage)
	}
}
