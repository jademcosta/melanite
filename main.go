package main

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jademcosta/melanite/converter"
	"github.com/julienschmidt/httprouter"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"
)

func main() {
	http.ListenAndServe(":8080", GetApp())
}

func GetApp() http.Handler {
	r := httprouter.New()
	r.GET("/*fileUri", FetcherFunc)

	n := negroni.New(negroni.NewRecovery())
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(r)

	return n
}

func FetcherFunc(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := removePrefixSlash(p.ByName("fileUri"))

	if isEmpty(url) || !isValidUrl(url) {
		http.Error(rw, "Invalid image url provided", http.StatusBadRequest)
		return
	}

	response, err := getImage(&url)
	if err != nil {
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	if externalImageNotFound(response) {
		http.Error(rw, "", http.StatusNotFound)
		return
	}

	defer response.Body.Close()

	imgAsBytes, err := decodeImageFromBody(&response.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// format := response.Header.Get("Content-Type")
	//
	// if len(format) == 0 {
	// 	http.Error(rw, "Content-Type not defined by upstream server", http.StatusInternalServerError)
	// 	return
	// }

	if output, ok := r.URL.Query()["o"]; ok && len(output) > 0 {
		outputFormat := output[0]
		if !converter.IsValidImageEncoding(outputFormat) {
			http.Error(rw, "Image conversion format not supported", http.StatusBadRequest)
			return
		}

		imgAsBytes, err = converter.Convert(*imgAsBytes, outputFormat)
		if err != nil {
			http.Error(rw, "Image conversion failed", http.StatusInternalServerError)
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
