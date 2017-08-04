package resizer_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/h2non/bimg"
	"github.com/jademcosta/melanite/resizer"
	"github.com/stretchr/testify/assert"
)

const testImagesFolder = "../test/images"

func TestItCanEnlargeImages(t *testing.T) {
	diskImage, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XS.png"))
	if err != nil {
		panic(err)
	}

	var imageEnlargementTests = []struct {
		newSize        string
		expectedWidth  int
		expectedHeight int
		messageForP    string
	}{
		{"520x0", 520, 294, "should keep aspect ratio when one of the sizes is zero"},
		{"0x294", 520, 294, "should keep aspect ratio when one of the sizes is zero"},
		{"520x200", 354, 200,
			"should use size with the smallest ratio as the leading one, and keep aspect ratio"},
		{"520x1500", 520, 294,
			"should use size with the smallest ratio as the leading one, and keep aspect ratio"},
	}
	//260 x 147
	for _, testCase := range imageEnlargementTests {

		resizedImage, err := resizer.Resize(diskImage, testCase.newSize)
		if err != nil {
			panic(err)
		}

		newSize, err := bimg.NewImage(resizedImage).Size()
		if err != nil {
			panic(err)
		}

		assert.Equal(t,
			testCase.expectedWidth, newSize.Width,
			fmt.Sprintf("The image should have a width of %d", testCase.expectedWidth))

		assert.Equal(t,
			testCase.expectedHeight, newSize.Height,
			fmt.Sprintf("The image should have a height of %d", testCase.expectedHeight))

		assert.Equal(t,
			"image/png", http.DetectContentType(resizedImage),
			"The image should not be converted to another format")
	}
}

func TestItCanDownsizeImages(t *testing.T) {
	//260 x 147
	diskImage, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-S.jpg"))
	if err != nil {
		panic(err)
	}

	var imageEnlargementTests = []struct {
		newSize        string
		expectedWidth  int
		expectedHeight int
		messageForP    string
	}{
		{"130x0", 130, 73, "should keep aspect ratio when one of the sizes is zero"},
		{"0x73", 130, 73, "should keep aspect ratio when one of the sizes is zero"},
		{"130x147", 130, 147,
			"should downsize keeping ratio and filling the remaining space with black background"},
		{"260x73", 260, 73,
			"should downsize keeping ratio and filling the remaining space with black background"},
	}

	for _, testCase := range imageEnlargementTests {

		resizedImage, err := resizer.Resize(diskImage, testCase.newSize)
		if err != nil {
			panic(err)
		}

		newSize, err := bimg.NewImage(resizedImage).Size()
		if err != nil {
			panic(err)
		}

		assert.Equal(t,
			testCase.expectedWidth, newSize.Width,
			fmt.Sprintf("The image should have a width of %d", testCase.expectedWidth))

		assert.Equal(t,
			testCase.expectedHeight, newSize.Height,
			fmt.Sprintf("The image should have a height of %d", testCase.expectedHeight))

		assert.Equal(t,
			"image/jpeg", http.DetectContentType(resizedImage),
			"The image should not be converted to another format")
	}
}

func TestItKeepSizesTheSameIfZeroIsProvided(t *testing.T) {
	//260 x 147
	diskImage, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XS.webp"))
	if err != nil {
		panic(err)
	}

	resizedImage, err := resizer.Resize(diskImage, "0x0")
	if err != nil {
		panic(err)
	}

	newSize, err := bimg.NewImage(resizedImage).Size()
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 260, newSize.Width, "The image should have a width of 260")

	assert.Equal(t, 147, newSize.Height, "The image should have a height of 147")

	assert.Equal(t,
		"image/webp", http.DetectContentType(resizedImage),
		"The image should not be converted to another format")
}

func TestItReturnsErrorIfWeUseIncorrectSizeFormat(t *testing.T) {
	diskImage, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", testImagesFolder, "park-view-XS.png"))
	if err != nil {
		panic(err)
	}

	var imageEnlargementSizeParams = []string{"520x", "x294", "0x", "x0", "x",
		"", "ax12", "1xb", "fdfxdfg", "xxxxxxxxx", "500xx600"}

	for _, testCase := range imageEnlargementSizeParams {
		_, err := resizer.Resize(diskImage, testCase)
		assert.NotNil(t, err,
			fmt.Sprintf("The input %s should return an error", testCase))
	}
}

func TestItReturnsErrorIfWeSendInvalidImage(t *testing.T) {
	_, err := resizer.Resize([]byte{23, 11, 32}, "23x44")
	assert.NotNil(t, err, "Should return an error due to invalid image")
}
