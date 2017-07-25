package main

import (
	"fmt"
	"net/http"

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
	filename := p.ByName("file")

	if filename == "/" {
		http.Error(rw, "No", http.StatusBadRequest)
		return
	}

	fmt.Fprintln(rw, fmt.Sprintf("Hello %s", filename))
}
