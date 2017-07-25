package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func main() {
	http.ListenAndServe(":8080", GetApp())
}

func GetApp() http.Handler {
	r := httprouter.New()
	r.GET("/*file", FetcherFunc)
	return r
}

func FetcherFunc(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	filename := removePrefixSlash(p.ByName("file"))

	if filename == "" {
		http.Error(rw, "No image url provided", http.StatusBadRequest)
		return
	}

	fmt.Fprintln(rw, fmt.Sprintf("Hello %s", filename))
}

func removePrefixSlash(s string) string {
	return strings.TrimPrefix(s, "/")
}
