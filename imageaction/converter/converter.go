package converter

import "github.com/h2non/bimg"

func IsValidImageEncoding(encoding string) bool {
	if encoding == "jpg" || encoding == "png" || encoding == "webp" {
		return true
	}
	return false
}

func Convert(imgAsBytes []byte, outputFormat string) (*[]byte, error) {

	var imgFormat bimg.ImageType = bimg.PNG

	switch outputFormat {
	case "png":
		imgFormat = bimg.PNG
	case "jpg":
		imgFormat = bimg.JPEG
	case "webp":
		imgFormat = bimg.WEBP
	}

	bts, err := bimg.Resize(imgAsBytes, bimg.Options{Type: imgFormat})
	if err != nil {
		return nil, err
	}

	return &bts, nil
}
