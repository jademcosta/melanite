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

		configuration, err := config.New(configContent, "")
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

		_, err = config.New(configContent, "")

		assert.Error(t, err, testCase.testMessage)
		assert.Equal(t, err.Error(),
			testCase.expectedErrorMessage, testCase.testMessage)
	}
}

func TestConfigInitializerOverridesImageSource(t *testing.T) {

	var imageSourceOnConfigTests = []struct {
		filename            string
		imgSourceFromArgs   string
		expectedImageSource string
		testMessage         string
	}{
		{"full_correct_config.yaml", "", "http://example.com",
			"the image source should be the the one on the file"},
		{"full_correct_config.yaml", "http://another.com", "http://another.com",
			"the image source should be the the one on the second parameter"},
		{"empty_config.yaml", "http://another.com", "http://another.com",
			"the image source should be the the one on the second parameter"},
	}

	for _, testCase := range imageSourceOnConfigTests {
		configContent, err :=
			ioutil.ReadFile(fmt.Sprintf("%s%s", configFixtureFilesFolder, testCase.filename))
		if err != nil {
			panic(err)
		}

		configuration, err := config.New(configContent, testCase.imgSourceFromArgs)
		if err != nil {
			panic(err)
		}

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

		configuration, err := config.New(configContent, "")
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

		_, err = config.New(configContent, "")

		assert.Error(t, err, testCase.testMessage)
		assert.Equal(t, err.Error(),
			testCase.expectedErrorMessage, testCase.testMessage)
	}
}
