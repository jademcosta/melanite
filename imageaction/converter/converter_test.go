package converter_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jademcosta/melanite/imageaction/converter"
	"github.com/stretchr/testify/assert"
)

const testImagesFolder = "../../test/images"

func TestValidationWorksForImageEncodingsThatAreSupported(t *testing.T) {

	var imageEncodingValidationTests = []struct {
		format         string
		expectedResult bool
	}{
		{"png", true},
		{"jpg", true},
		{"", false},
		{"jpeg", false},
		{"gif", false},
		{"bmp", false},
		{"svg", false},
	}

	for _, sample := range imageEncodingValidationTests {
		assert.Equal(t,
			sample.expectedResult, converter.IsValidImageEncoding(sample.format),
			fmt.Sprintf("should be %v", sample.expectedResult))
	}
}

func TestConversionFromJpgToPng(t *testing.T) {
	diskImage, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XXS.jpg"))
	if err != nil {
		panic(err)
	}

	convertedImage, err := converter.Convert(diskImage, "png")
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"image/png", http.DetectContentType(*convertedImage),
		"The image should be converted to PNG")
}

func TestConversionFromPngToJpg(t *testing.T) {
	diskImage, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XXS.png"))
	if err != nil {
		panic(err)
	}

	convertedImage, err := converter.Convert(diskImage, "jpg")
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"image/jpeg", http.DetectContentType(*convertedImage),
		"The image should be converted to JPG")
}

func TestConversionFromPngToWebp(t *testing.T) {
	diskImage, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XXS.png"))
	if err != nil {
		panic(err)
	}

	convertedImage, err := converter.Convert(diskImage, "webp")
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"image/webp", http.DetectContentType(*convertedImage),
		"The image should be converted to WEBP")
}

func TestConversionFromJpgToWebp(t *testing.T) {
	diskImage, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XXS.jpg"))
	if err != nil {
		panic(err)
	}

	convertedImage, err := converter.Convert(diskImage, "webp")
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"image/webp", http.DetectContentType(*convertedImage),
		"The image should be converted to WEBP")
}

func TestConversionFromWebpToJpg(t *testing.T) {
	diskImage, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XXS.webp"))
	if err != nil {
		panic(err)
	}

	convertedImage, err := converter.Convert(diskImage, "jpg")
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"image/jpeg", http.DetectContentType(*convertedImage),
		"The image should be converted to JPG")
}

func TestConversionFromWebpToPng(t *testing.T) {
	diskImage, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XXS.webp"))
	if err != nil {
		panic(err)
	}

	convertedImage, err := converter.Convert(diskImage, "png")
	if err != nil {
		panic(err)
	}

	assert.Equal(t,
		"image/png", http.DetectContentType(*convertedImage),
		"The image should be converted to PNG")
}

func TestConversionOfInvalidFileReturnError(t *testing.T) {
	wrongImage := []byte{12, 34, 124}

	_, err := converter.Convert(wrongImage, "jpg")
	assert.NotNil(t, err, "An error should be returned")

	_, err = converter.Convert(wrongImage, "png")
	assert.NotNil(t, err, "An error should be returned")
}
