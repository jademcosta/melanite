package converter

import (
	"fmt"

	"github.com/h2non/bimg"
	"github.com/jademcosta/melanite/imageaction"
)

func Convert(imgAsBytes []byte, outputFormat string) (*[]byte, error) {

	var imgFormat bimg.ImageType = bimg.PNG

	switch outputFormat {
	case "png":
		imgFormat = bimg.PNG
	case "jpg":
		imgFormat = bimg.JPEG
	case "webp":
		imgFormat = bimg.WEBP
	default:
		return nil, imageaction.Error{
			Message: fmt.Sprintf("converter: invalid output format: %s", outputFormat)}
	}

	bts, err := bimg.Resize(imgAsBytes, bimg.Options{Type: imgFormat})
	if err != nil {
		return nil, err
	}

	return &bts, nil
}
