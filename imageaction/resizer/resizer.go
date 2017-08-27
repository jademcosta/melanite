package resizer

import (
	"strconv"
	"strings"

	"github.com/h2non/bimg"
)

func Resize(imgAsBytes []byte, newSize string) ([]byte, error) {

	sizes := strings.Split(newSize, "x")

	width, err := strconv.Atoi(sizes[0])
	if err != nil {
		return nil, err
	}

	height, err := strconv.Atoi(sizes[1])
	if err != nil {
		return nil, err
	}

	currentSize, err := bimg.NewImage(imgAsBytes).Size()
	if err != nil {
		return nil, err
	}

	currentWidth := currentSize.Width
	currentHeight := currentSize.Height

	image := bimg.NewImage(imgAsBytes)

	var resizedImage []byte

	if (width > currentWidth || width == 0) && (height > currentHeight || height == 0) {
		resizedImage, err = image.Enlarge(width, height)
	} else {
		resizedImage, err = image.Resize(width, height)
	}

	if err != nil {
		return nil, err
	}

	return resizedImage, nil
}
