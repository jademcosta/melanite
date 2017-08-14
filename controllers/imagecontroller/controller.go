package imagecontroller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jademcosta/melanite/config"
	"github.com/jademcosta/melanite/imageaction/converter"
	"github.com/jademcosta/melanite/imageaction/resizer"
	log "github.com/sirupsen/logrus"
)

type ImageController struct {
	config config.Config
	logger *log.Logger
}

func New(config config.Config, logger *log.Logger) *ImageController {
	return &ImageController{config: config, logger: logger}
}

func (controller *ImageController) ServeHTTP(rw http.ResponseWriter,
	r *http.Request) {
	logger := controller.logger

	filePath := r.URL.Path
	emptyFilePath := "/"
	if filePath == emptyFilePath {
		logger.Info("Empty image path.")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s%s", controller.config.ImageSource, filePath)

	response, err := getImage(&url)
	if err != nil {
		logger.Errorf("Error when trying to get image from %s. Error: %s", url, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response.StatusCode != http.StatusOK {
		logger.Info("Image at %s answered %d", url, response.StatusCode)
		rw.WriteHeader(response.StatusCode)
		return
	}

	defer response.Body.Close()

	imgAsBytes, err := decodeImageFromBody(&response.Body)
	if err != nil {
		logger.Errorf("Error when trying to decode image. Error: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if output, ok := r.URL.Query()["o"]; ok && len(output) > 0 {
		outputFormat := output[0]
		if !converter.IsValidImageEncoding(outputFormat) {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		imgAsBytes, err = converter.Convert(*imgAsBytes, outputFormat)
		if err != nil {
			logger.Errorf("Error when trying to convert image. Error: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if resizeParam, ok := r.URL.Query()["r"]; ok && len(resizeParam) > 0 {
		resizeDimensions := resizeParam[0]

		*imgAsBytes, err = resizer.Resize(*imgAsBytes, resizeDimensions)
		if err != nil {
			logger.Errorf("Error when trying to resize image. Error: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	rw.Header().Add("Content-Length", strconv.Itoa(len(*imgAsBytes)))
	rw.Header().Add("Content-Type", http.DetectContentType(*imgAsBytes))
	rw.Write(*imgAsBytes)
}

func getImage(url *string) (*http.Response, error) {
	return http.Get(*url)
}

func decodeImageFromBody(body *io.ReadCloser) (*[]byte, error) {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(*body)

	if err != nil {
		return nil, err
	}

	b := buf.Bytes()
	return &b, nil
}
