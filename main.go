package main

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jademcosta/melanite/converter"
	"github.com/jademcosta/melanite/resizer"
	"github.com/julienschmidt/httprouter"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const defaultLogLevel = log.InfoLevel

var defaultLogFormatter = &log.JSONFormatter{}

func main() {
	http.ListenAndServe(":8080", GetApp(defaultLogLevel, defaultLogFormatter))
}

func GetApp(logLevel log.Level, logFormatter log.Formatter) http.Handler {
	r := httprouter.New()
	r.GET("/*fileUri", FetcherFunc)

	n := negroni.New(negroni.NewRecovery())
	appLog := log.New()
	appLog.SetLevel(logLevel)
	appLog.Formatter = logFormatter

	n.Use(negronilogrus.NewMiddlewareFromLogger(appLog, "melanite"))
	n.UseHandler(r)

	return n
}

func FetcherFunc(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := removePrefixSlash(p.ByName("fileUri"))

	if isEmpty(url) || !isValidUrl(url) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

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

func getImage(url *string) (*http.Response, error) {
	return http.Get(*url)
}

func isEmpty(s string) bool {
	if s == "" {
		return true
	}
	return false
}

func isValidUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return true
}

func removePrefixSlash(s string) string {
	return strings.TrimPrefix(s, "/")
}
