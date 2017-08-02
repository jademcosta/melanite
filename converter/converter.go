package converter

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
)

func IsValidImageEncoding(encoding string) bool {
	if encoding == "jpg" || encoding == "png" {
		return true
	}
	return false
}

func Convert(imgAsBytes []byte, outputFormat string) (*[]byte, error) {

	image, _, err := image.Decode(bytes.NewReader(imgAsBytes))
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	switch outputFormat {
	case "png":
		err = png.Encode(buf, image)
	case "jpg":
		err = jpeg.Encode(buf, image, nil)
	}

	if err != nil {
		return nil, err
	}

	convertedImg := buf.Bytes()

	return &convertedImg, nil
}
