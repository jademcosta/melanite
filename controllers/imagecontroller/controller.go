package imagecontroller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/jademcosta/melanite/config"
	"github.com/jademcosta/melanite/imageaction"
	"github.com/jademcosta/melanite/imageaction/converter"
	"github.com/jademcosta/melanite/imageaction/resizer"
	log "github.com/sirupsen/logrus"
)

const urlQueryParamKeyOutput string = "out"
const urlQueryParamKeyResize string = "res"

const responseHeaderContentLength string = "Content-Length"
const responseHeaderContentType string = "Content-Type"

const imageRequestTimeout time.Duration = 30 * time.Second

var httpClient http.Client

type ImageController struct {
	config config.Config
	logger *log.Logger
}

func init() {
	httpClient = http.Client{
		Timeout: time.Duration(imageRequestTimeout),
	}
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

	if output, ok := r.URL.Query()[urlQueryParamKeyOutput]; ok && len(output) > 0 {
		outputFormat := output[0]

		imgAsBytes, err = converter.Convert(*imgAsBytes, outputFormat)
		if err != nil {
			logger.Errorf("Error when trying to convert image. Error: %s", err)
			switch err.(type) {
			case imageaction.Error:
				rw.WriteHeader(http.StatusBadRequest)
			default:
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}

	if resizeParam, ok := r.URL.Query()[urlQueryParamKeyResize]; ok && len(resizeParam) > 0 {
		resizeDimensions := resizeParam[0]

		*imgAsBytes, err = resizer.Resize(*imgAsBytes, resizeDimensions)
		if err != nil {
			logger.Errorf("Error when trying to resize image. Error: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	rw.Header().Add(responseHeaderContentLength, strconv.Itoa(len(*imgAsBytes)))
	rw.Header().Add(responseHeaderContentType, http.DetectContentType(*imgAsBytes))
	rw.Write(*imgAsBytes)
}

func getImage(url *string) (*http.Response, error) {
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return nil, err
	}
	// Avoid connections hanging open
	req.Header.Set("Connection", "close")

	return httpClient.Do(req)
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
