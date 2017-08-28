package config_test

import (
	"fmt"
	"io/ioutil"
	"os"
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

func TestConfigInitializerForValidImageSourceKeyFromEnvVar(t *testing.T) {

	var imageSourceOnConfigTests = []struct {
		filename            string
		envVarValue         string
		expectedImageSource string
		testMessage         string
	}{
		{"full_correct_config.yaml", "http://other.com", "http://other.com",
			"the image source on env var should replace the one in config file"},
		{"incorrect_image_source_config.yaml", "http://other.com",
			"http://other.com",
			"the image source on env var should replace the one in config file"},
		{"empty_config.yaml", "http://other.com", "http://other.com",
			"the image source on env var should replace the one in config file"},
	}

	for _, testCase := range imageSourceOnConfigTests {
		configContent, err :=
			ioutil.ReadFile(fmt.Sprintf("%s%s", configFixtureFilesFolder, testCase.filename))
		if err != nil {
			panic(err)
		}

		os.Setenv(config.EnvVarKeyImageSource, testCase.envVarValue)

		configuration, err := config.New(configContent)
		if err != nil {
			panic(err)
		}
		os.Unsetenv(config.EnvVarKeyImageSource)
		assert.Equal(t, testCase.expectedImageSource,
			configuration.ImageSource, testCase.testMessage)
	}
}

func TestConfigInitializerForValidPortKey(t *testing.T) {

	var imageSourceOnConfigTests = []struct {
		filename     string
		expectedPort string
		testMessage  string
	}{
		{"full_correct_config.yaml", "80", "the port should be the equal"},
		{"empty_port_config.yaml", "", "the port should be the equal"},
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

		assert.Equal(t, testCase.expectedPort,
			configuration.Port, testCase.testMessage)
	}
}

func TestConfigInitializerForInvalidPortKey(t *testing.T) {

	var imageSourceOnConfigTests = []struct {
		filename             string
		expectedErrorMessage string
		testMessage          string
	}{
		{"incorrect_port_config.yaml",
			"config: port should be an integer",
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

func TestAnEmptyConfigFileWorksIfImageSourceEnvVarIsSet(t *testing.T) {

	os.Setenv(config.EnvVarKeyImageSource, "http://site.com")

	configuration, err := config.New([]byte{})

	os.Unsetenv(config.EnvVarKeyImageSource)
	assert.Nil(t, err,
		"There should be no error, as the image source was set on env var")
	assert.Equal(t, "http://site.com",
		configuration.ImageSource,
		"the image source on env var should compensate the abscence of config file")

}
