package imagecontroller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jademcosta/melanite/config"
	"github.com/jademcosta/melanite/converter"
	"github.com/jademcosta/melanite/resizer"
	"github.com/julienschmidt/httprouter"
)

type ImageController struct {
	config config.Config
}

func New(config config.Config) *ImageController {
	return &ImageController{config: config}
}

func (controller *ImageController) ServeHttp(rw http.ResponseWriter,
	r *http.Request, p httprouter.Params) {

	filePath := p.ByName("fileUri")
	if filePath == "/" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s%s", controller.config.ImageSource, filePath)

	response, err := getImage(&url)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if externalImageNotFound(response) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	defer response.Body.Close()

	imgAsBytes, err := decodeImageFromBody(&response.Body)
	if err != nil {
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
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if resizeParam, ok := r.URL.Query()["r"]; ok && len(resizeParam) > 0 {
		resizeDimensions := resizeParam[0]

		*imgAsBytes, err = resizer.Resize(*imgAsBytes, resizeDimensions)
		if err != nil {
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

func externalImageNotFound(response *http.Response) bool {
	return response.StatusCode == http.StatusNotFound
}
