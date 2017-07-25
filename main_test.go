package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

// func TestMain(m *testing.M) {
//     a = main.App{}
//     a.Initialize(
//         os.Getenv("TEST_DB_USERNAME"),
//         os.Getenv("TEST_DB_PASSWORD"),
//         os.Getenv("TEST_DB_NAME"))
//
//     ensureTableExists()
//
//     code := m.Run()
//
//     clearTable()
//
//     os.Exit(codunc TestMain(m *testing.M) {
//     a = main.App{}
//     a.Initialize(
//         os.Getenv("TEST_DB_USERNAME"),
//         os.Getenv("TEST_DB_PASSWORD"),
//         os.Getenv("TEST_DB_NAME"))
//
//     ensureTableExists()
//
//     code := m.Run()
//
//     clearTable()
//
//     os.Exit(code)
// }e)
// }

func TestSimple(t *testing.T) {
	req, _ := http.NewRequest("GET", "/doesntmatter:(", nil)
	response := httptest.NewRecorder()

	oi := []httprouter.Param{httprouter.Param{Key: "file", Value: "jade"}}

	FetcherFunc(response, req, oi)

	assert.Equal(t, response.Body.String(), "Hello jade\n", "they should be equal")
}

func TestFromOutside(t *testing.T) {
	ts := httptest.NewServer(GetApp())
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/%s", ts.URL, "jade"))
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 200, res.StatusCode, "status code should be 200")
	assert.Equal(t, string(greeting), "Hello jade\n", "they should be equal")
}

func TestNoUrlReturnsBadRequest(t *testing.T) {
	ts := httptest.NewServer(GetApp())
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/%s", ts.URL, ""))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 400, res.StatusCode, "status code should be 400")
	assert.Equal(t, string(body), "No image url provided\n", "they should be equal")
}
