package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func main() {
	http.ListenAndServe(":8080", GetApp())
}

func GetApp() http.Handler {
	r := httprouter.New()
	r.GET("/*fileUri", FetcherFunc)
	return r
}

func FetcherFunc(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := removePrefixSlash(p.ByName("fileUri"))

	if isEmpty(url) {
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
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	format := response.Header.Get("Content-type")

	if len(format) == 0 {
		http.Error(rw, "Content-type not defined by upstream server", http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-type", format)
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

func removePrefixSlash(s string) string {
	return strings.TrimPrefix(s, "/")
}
